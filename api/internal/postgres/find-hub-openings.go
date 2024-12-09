package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) FindHubOpenings(
	ctx context.Context,
	req *hub.FindHubOpeningsRequest,
) ([]hub.HubOpening, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("hub user not found in context")
		return nil, db.ErrInternal
	}

	query := `
SELECT
	o.id as opening_id_within_company,
	d.domain_name as company_domain,
	e.company_name as company_name,
	o.title as job_title,
	o.jd as jd,
	o.pagination_key
FROM openings o
	JOIN hub_users hu ON hu.id = $1
	JOIN employers e ON o.employer_id = e.id
	JOIN employer_primary_domains epd ON e.id = epd.employer_id
	JOIN domains d ON epd.domain_id = d.id
	JOIN opening_locations ol ON o.employer_id = ol.employer_id AND o.id = ol.opening_id
	JOIN locations l ON ol.location_id = l.id
	WHERE o.state = 'ACTIVE_OPENING_STATE'
		AND l.country_code = COALESCE($2, hu.resident_country_code)
	`

	// $1 is hub_users.id (needed for getting the resident_country_code and
	// the resident_city of the logged in Hub User)
	args := []interface{}{hubUser.ID}

	// $2 is for country_code
	if req.CountryCode != nil {
		args = append(args, string(*req.CountryCode))
	} else {
		// If no country code is passed, we will fallback to the Hub User's resident country
		args = append(args, nil)
	}

	// $1 and $2 already filled, so we start from $3
	argPos := 3

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
			req.SalaryRange.Min,
			req.SalaryRange.Max,
		)
		argPos += 3
		p.log.Dbg("salary", "salaryConds", salaryConds)
	}
	p.log.Dbg("salary", "query", query, "args", args, "argPos", argPos)

	if len(req.RemoteCountryCodes) > 0 {
		// TODO: Accomodate Globally Remote Openings
		placeholders := make([]string, len(req.RemoteCountryCodes))
		for i := range req.RemoteCountryCodes {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.RemoteCountryCodes[i])
			argPos++
		}
		remoteCountryConds := fmt.Sprintf(
			"o.remote_country_codes && ARRAY[%s]::text[]",
			strings.Join(placeholders, ","),
		)
		whereConditions = append(whereConditions, remoteCountryConds)
		p.log.Dbg("remote_country_codes", "Conds", remoteCountryConds)
	}
	p.log.Dbg("remote_country", "query", query, "args", args, "argPos", argPos)

	if len(req.RemoteTimezones) > 0 {
		// TODO: Accomodate Globally Remote Openings
		placeholders := make([]string, len(req.RemoteTimezones))
		for i := range req.RemoteTimezones {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.RemoteTimezones[i])
			argPos++
		}
		remoteTimezoneConds := fmt.Sprintf(
			"o.remote_timezones && ARRAY[%s]::text[]",
			strings.Join(placeholders, ","),
		)
		whereConditions = append(whereConditions, remoteTimezoneConds)
		p.log.Dbg("remote_timezone", "remoteTimezoneConds", remoteTimezoneConds)
	}
	p.log.Dbg("remote_timezone", "query", query, "args", args, "argPos", argPos)

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

	var openings []hub.HubOpening
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
