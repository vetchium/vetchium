package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) FindHubOpenings(
	ctx context.Context,
	req *hub.FindHubOpeningsRequest,
) ([]hub.HubOpening, error) {
	query := `
WITH applicable_openings AS (
	SELECT
		o.id as opening_id,
		o.employer_id,
		can_apply($1, o.employer_id, o.id) as can_apply
	FROM openings o
	WHERE o.state = $2
)
SELECT
	o.id as opening_id_within_company,
	d.domain_name as company_domain,
	e.company_name as company_name,
	o.title as job_title,
	o.jd as jd,
	o.pagination_key
FROM openings o
	JOIN employers e ON o.employer_id = e.id
	JOIN employer_primary_domains epd ON e.id = epd.employer_id
	JOIN domains d ON epd.domain_id = d.id
	JOIN applicable_openings ao ON o.id = ao.opening_id AND o.employer_id = ao.employer_id
	LEFT JOIN opening_locations ol ON o.employer_id = ol.employer_id AND o.id = ol.opening_id
	LEFT JOIN locations l ON ol.location_id = l.id
	LEFT JOIN opening_tag_mappings otm ON o.employer_id = otm.employer_id AND o.id = otm.opening_id
	LEFT JOIN opening_tags ot ON otm.tag_id = ot.id
	WHERE o.state = $2
		AND ao.can_apply = true
		AND (
			COALESCE(l.country_code, '') = $3
			OR $4 = ANY(o.remote_country_codes)
			OR $3 = ANY(o.remote_country_codes)
		)`

	args := []interface{}{}

	// $1 is for hub_user_id from context
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return nil, db.ErrInternal
	}
	args = append(args, hubUser.ID)

	// $2 is for opening state
	args = append(args, common.ActiveOpening)

	// $3 is for country_code
	args = append(args, string(req.CountryCode))

	// $4 is for global country code
	args = append(args, string(common.GlobalCountryCode))

	// Previous parameters filled, so we start from $5
	argPos := 5

	// Add city filter if specified
	if len(req.Cities) > 0 {
		cityConditions := make([]string, len(req.Cities))
		for i, city := range req.Cities {
			// TODO: Check if doing this condition in reverse is better, like we do for company_domains
			// TODO: We would need some kind of validation or punishment for people who abuse the city_aka field to grab more eyeballs
			cityConditions[i] = fmt.Sprintf(
				"$%d = ANY(l.city_aka)",
				argPos,
			)
			args = append(args, city)
			argPos++
		}
		query += " AND (" + strings.Join(cityConditions, " OR ") + ")"
	}
	p.log.Dbg("cityfilter", "query", query, "args", args, "argPos", argPos)

	// Add other filters
	whereConditions := []string{}

	// Add tag filter if specified
	if len(req.Tags) > 0 {
		placeholders := make([]string, len(req.Tags))
		for i := range req.Tags {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.Tags[i])
			argPos++
		}
		// Use EXISTS to properly handle the LEFT JOIN case and match against tag IDs
		tagConds := fmt.Sprintf(
			"EXISTS (SELECT 1 FROM opening_tag_mappings otm2 WHERE otm2.employer_id = o.employer_id AND otm2.opening_id = o.id AND otm2.tag_id = ANY(ARRAY[%s]::uuid[]))",
			strings.Join(placeholders, ","),
		)
		whereConditions = append(whereConditions, tagConds)
		p.log.Dbg("tag conditions", "tagConds", tagConds)
	}

	// Add term filter if specified
	if len(req.Terms) > 0 {
		termConditions := make([]string, len(req.Terms))
		for i, term := range req.Terms {
			termConditions[i] = fmt.Sprintf(
				"o.title ILIKE $%d",
				argPos,
			)
			args = append(args, "%"+term+"%")
			argPos++
		}
		termConds := "(" + strings.Join(termConditions, " OR ") + ")"
		whereConditions = append(whereConditions, termConds)
		p.log.Dbg("term conditions", "termConds", termConds)
	}

	// If we have both tags and terms, we want to match either of them
	if len(req.Tags) > 0 && len(req.Terms) > 0 {
		// Take the last two conditions (tags and terms) and combine them with OR
		lastIdx := len(whereConditions) - 1
		tagCond := whereConditions[lastIdx-1]
		termCond := whereConditions[lastIdx]
		whereConditions = whereConditions[:lastIdx-1]
		whereConditions = append(
			whereConditions,
			fmt.Sprintf("(%s OR %s)", tagCond, termCond),
		)
	}

	if len(req.OpeningTypes) > 0 {
		placeholders := make([]string, len(req.OpeningTypes))
		for i := range req.OpeningTypes {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.OpeningTypes[i])
			argPos++
		}
		opTypeConds := fmt.Sprintf(
			"o.opening_type = ANY(ARRAY[%s]::opening_types[])",
			strings.Join(placeholders, ","),
		)
		whereConditions = append(whereConditions, opTypeConds)
		p.log.Dbg("opening_type conditions", "opTypeConds", opTypeConds)
	}
	p.log.Dbg("opening_type", "query", query, "args", args, "argPos", argPos)

	if len(req.CompanyDomains) > 0 {
		placeholders := make([]string, len(req.CompanyDomains))
		for i := range req.CompanyDomains {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.CompanyDomains[i])
			argPos++
		}
		domainConds := fmt.Sprintf(
			"d.domain_name = ANY(ARRAY[%s])",
			strings.Join(placeholders, ","),
		)
		whereConditions = append(whereConditions, domainConds)
		p.log.Dbg("company_domain", "domainConds", domainConds)
	}
	p.log.Dbg("company_domain", "query", query, "args", args, "argPos", argPos)

	if req.ExperienceRange != nil {
		expConds := fmt.Sprintf(
			"o.yoe_min <= $%d AND o.yoe_max >= $%d",
			argPos+1,
			argPos,
		)
		whereConditions = append(whereConditions, expConds)
		args = append(
			args,
			req.ExperienceRange.YoeMin,
			req.ExperienceRange.YoeMax,
		)
		argPos += 2
		p.log.Dbg("experience", "expConds", expConds)
	}
	p.log.Dbg("experience", "query", query, "args", args, "argPos", argPos)

	if req.SalaryRange != nil {
		salaryConds := fmt.Sprintf(
			"o.salary_currency = $%d AND o.salary_min >= $%d AND o.salary_max <= $%d",
			argPos,
			argPos+1,
			argPos+2,
		)
		whereConditions = append(whereConditions, salaryConds)
		args = append(
			args,
			req.SalaryRange.Currency,
			req.SalaryRange.MinAmount,
			req.SalaryRange.MaxAmount,
		)
		argPos += 3
		p.log.Dbg("salary", "salaryConds", salaryConds)
	}
	p.log.Dbg("salary", "query", query, "args", args, "argPos", argPos)

	if req.MinEducationLevel != nil {
		minEduConds := fmt.Sprintf("o.min_education_level = $%d", argPos)
		whereConditions = append(whereConditions, minEduConds)
		args = append(args, *req.MinEducationLevel)
		argPos++
		p.log.Dbg("min_education_level", "minEduConds", minEduConds)
	}
	p.log.Dbg("min_education", "query", query, "args", args, "argPos", argPos)
	// End of all the Where clauses. Now we need to append them to the query

	if len(whereConditions) > 0 {
		query += " AND " + strings.Join(whereConditions, " AND ")
	}
	p.log.Dbg("with WHERE", "query", query, "args", args, "argPos", argPos)

	// Add pagination and ordering
	query += fmt.Sprintf(" AND o.pagination_key > $%d", argPos)
	args = append(args, req.PaginationKey)
	argPos++

	// Add GROUP BY
	query += `
		GROUP BY
			o.employer_id,
			o.id,
			o.title,
			o.jd,
			d.domain_name,
			e.company_name,
			o.pagination_key
		ORDER BY o.pagination_key
`

	// Add LIMIT
	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	p.log.Dbg("Final hub openings query", "query", query, "args", args)

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("error querying openings", "err", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	openings := []hub.HubOpening{}
	for rows.Next() {
		var opening hub.HubOpening
		err := rows.Scan(
			&opening.OpeningIDWithinCompany,
			&opening.CompanyDomain,
			&opening.CompanyName,
			&opening.JobTitle,
			&opening.JD,
			&opening.PaginationKey,
		)
		if err != nil {
			p.log.Err("error scanning opening row", "err", err)
			return nil, db.ErrInternal
		}
		openings = append(openings, opening)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating opening rows", "err", err)
		return nil, db.ErrInternal
	}

	return openings, nil
}
