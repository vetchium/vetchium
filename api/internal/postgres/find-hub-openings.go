package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) FindHubOpenings(
	ctx context.Context,
	req *vetchi.FindHubOpeningsRequest,
) ([]vetchi.HubOpening, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("hub user not found in context")
		return nil, db.ErrInternal
	}

	query := `
SELECT
	o.id as opening_id_within_company,
	epd.domain_name as company_domain,
	e.company_name as company_name,
	o.title as job_title,
	o.jd as jd,
	o.pagination_key
FROM openings o
	JOIN hub_users hu ON hu.id = $1
	JOIN employers e ON o.employer_id = e.id
	JOIN employer_primary_domains epd ON e.id = epd.employer_id
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
			// TODO: Check if doing this condition in reverse is better
			cityConditions[i] = fmt.Sprintf(
				"$%d = ANY(l.city_aka)",
				argPos,
			)
			args = append(args, city)
			argPos++
		}
		query += " AND (" + strings.Join(cityConditions, " OR ") + ")"
	}
	// Now after exiting the above loop, argPos will be argPos + len(req.Cities)

	// Add other filters
	whereConditions := []string{}

	if len(req.OpeningTypes) > 0 {
		placeholders := make([]string, len(req.OpeningTypes))
		for i := range req.OpeningTypes {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.OpeningTypes[i])
			argPos++
		}
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.opening_type = ANY(ARRAY[%s]::opening_types[])",
				strings.Join(placeholders, ","),
			),
		)
	}

	if len(req.CompanyDomains) > 0 {
		placeholders := make([]string, len(req.CompanyDomains))
		for i := range req.CompanyDomains {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.CompanyDomains[i])
			argPos++
		}
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"d.domain_name = ANY(ARRAY[%s])",
				strings.Join(placeholders, ","),
			),
		)
	}

	if req.ExperienceRange != nil {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.yoe_min <= $%d AND o.yoe_max >= $%d",
				argPos+1,
				argPos,
			),
		)
		args = append(
			args,
			req.ExperienceRange.YoeMin,
			req.ExperienceRange.YoeMax,
		)
		argPos += 2
	}

	if req.SalaryRange != nil {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.salary_currency = $%d AND o.salary_min >= $%d AND o.salary_max <= $%d",
				argPos,
				argPos+1,
				argPos+2,
			),
		)
		args = append(
			args,
			req.SalaryRange.Currency,
			req.SalaryRange.Min,
			req.SalaryRange.Max,
		)
		argPos += 3
	}

	if len(req.RemoteCountryCodes) > 0 {
		placeholders := make([]string, len(req.RemoteCountryCodes))
		for i := range req.RemoteCountryCodes {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.RemoteCountryCodes[i])
			argPos++
		}
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.remote_country_codes && ARRAY[%s]::text[]",
				strings.Join(placeholders, ","),
			),
		)
	}

	if len(req.RemoteTimezones) > 0 {
		placeholders := make([]string, len(req.RemoteTimezones))
		for i := range req.RemoteTimezones {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, req.RemoteTimezones[i])
			argPos++
		}
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.remote_timezones && ARRAY[%s]::text[]",
				strings.Join(placeholders, ","),
			),
		)
	}

	if req.MinEducationLevel != nil {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf("o.min_education_level = $%d", argPos),
		)
		args = append(args, *req.MinEducationLevel)
		argPos++
	}

	if len(whereConditions) > 0 {
		query += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Add pagination and ordering
	query += fmt.Sprintf(" AND o.pagination_key > $%d", argPos)
	args = append(args, req.PaginationKey)
	argPos++

	// Add GROUP BY
	query += `
		GROUP BY o.id, o.title, d.domain_name, e.onboard_admin_email, o.pagination_key
		ORDER BY o.pagination_key
	`

	// Add LIMIT
	query += fmt.Sprintf(" LIMIT $%d", argPos)
	args = append(args, req.Limit)

	p.log.Dbg(
		"find hub openings query",
		"query",
		query,
		"args",
		args,
		"conditions",
		whereConditions,
	)

	// Add this query before applying filters to see what data exists
	debugQuery := `
		SELECT DISTINCT
			o.id,
			o.title,
			o.yoe_min,
			o.yoe_max,
			o.salary_currency,
			o.salary_min,
			o.salary_max,
			o.min_education_level,
			o.remote_country_codes,
			o.remote_timezones
		FROM openings o
		WHERE o.state = 'ACTIVE_OPENING_STATE'
	`
	debugRows, err := p.pool.Query(ctx, debugQuery)
	if err == nil {
		defer debugRows.Close()
		p.log.Dbg("existing openings in database:")
		for debugRows.Next() {
			var id, title string
			var yoeMin, yoeMax, salaryMin, salaryMax int
			var salaryCurrency, minEducation string
			var remoteCountries, remoteTimezones []string
			err := debugRows.Scan(
				&id,
				&title,
				&yoeMin,
				&yoeMax,
				&salaryCurrency,
				&salaryMin,
				&salaryMax,
				&minEducation,
				&remoteCountries,
				&remoteTimezones,
			)
			if err == nil {
				p.log.Dbg(
					"opening",
					"id",
					id,
					"title",
					title,
					"yoe",
					fmt.Sprintf("%d-%d", yoeMin, yoeMax),
					"salary",
					fmt.Sprintf(
						"%s %d-%d",
						salaryCurrency,
						salaryMin,
						salaryMax,
					),
					"education",
					minEducation,
					"countries",
					remoteCountries,
					"timezones",
					remoteTimezones,
				)
			}
		}
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying openings: %w", err)
	}
	defer rows.Close()

	var openings []vetchi.HubOpening
	for rows.Next() {
		var opening vetchi.HubOpening
		err := rows.Scan(
			&opening.OpeningIDWithinCompany,
			&opening.JobTitle,
			&opening.CompanyDomain,
			&opening.CompanyName,
			&opening.PaginationKey,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning opening row: %w", err)
		}
		openings = append(openings, opening)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating opening rows: %w", err)
	}

	return openings, nil
}
