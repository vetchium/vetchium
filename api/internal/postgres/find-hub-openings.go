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
	// Get hub user from context
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		return nil, fmt.Errorf("hub user not found in context")
	}

	query := `
		SELECT DISTINCT
			o.id as opening_id,
			o.title as job_title,
			d.domain_name as company_domain,
			e.onboard_admin_email as company_name,
			array_agg(DISTINCT l.title) as cities,
			o.pagination_key
		FROM openings o
		JOIN employers e ON o.employer_id = e.id
		JOIN domains d ON e.id = d.employer_id
		LEFT JOIN opening_locations ol ON o.employer_id = ol.employer_id AND o.id = ol.opening_id
		LEFT JOIN locations l ON ol.location_id = l.id
		JOIN hub_users hu ON hu.id = $1
		WHERE o.state = 'ACTIVE_OPENING_STATE'
		AND l.country_code = COALESCE($2, hu.resident_country_code)
	`

	args := []interface{}{hubUser.ID}
	if req.CountryCode != nil {
		args = append(args, string(*req.CountryCode))
	} else {
		args = append(args, nil)
	}
	argCount := 3

	// Add city filter if specified
	if len(req.Cities) > 0 {
		cityConditions := make([]string, len(req.Cities))
		for i, city := range req.Cities {
			cityConditions[i] = fmt.Sprintf(
				"$%d = ANY(l.city_aka)",
				argCount,
			)
			args = append(args, city)
			argCount++
		}
		query += " AND (" + strings.Join(cityConditions, " OR ") + ")"
	}

	// Add other filters
	whereConditions := []string{}

	if len(req.OpeningTypes) > 0 {
		placeholders := make([]string, len(req.OpeningTypes))
		for i := range req.OpeningTypes {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, req.OpeningTypes[i])
			argCount++
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
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, req.CompanyDomains[i])
			argCount++
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
				argCount+1,
				argCount,
			),
		)
		args = append(
			args,
			req.ExperienceRange.YoeMin,
			req.ExperienceRange.YoeMax,
		)
		argCount += 2
	}

	if req.SalaryRange != nil {
		whereConditions = append(
			whereConditions,
			fmt.Sprintf(
				"o.salary_currency = $%d AND o.salary_min >= $%d AND o.salary_max <= $%d",
				argCount,
				argCount+1,
				argCount+2,
			),
		)
		args = append(
			args,
			req.SalaryRange.Currency,
			req.SalaryRange.Min,
			req.SalaryRange.Max,
		)
		argCount += 3
	}

	if len(req.RemoteCountryCodes) > 0 {
		placeholders := make([]string, len(req.RemoteCountryCodes))
		for i := range req.RemoteCountryCodes {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, req.RemoteCountryCodes[i])
			argCount++
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
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, req.RemoteTimezones[i])
			argCount++
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
			fmt.Sprintf("o.min_education_level = $%d", argCount),
		)
		args = append(args, *req.MinEducationLevel)
		argCount++
	}

	if len(whereConditions) > 0 {
		query += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Add pagination and ordering
	query += fmt.Sprintf(" AND o.pagination_key > $%d", argCount)
	args = append(args, req.PaginationKey)
	argCount++

	// Add GROUP BY
	query += `
		GROUP BY o.id, o.title, d.domain_name, e.onboard_admin_email, o.pagination_key
		ORDER BY o.pagination_key
	`

	// Add LIMIT
	query += fmt.Sprintf(" LIMIT $%d", argCount)
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
		var cities []string
		err := rows.Scan(
			&opening.OpeningIDWithinCompany,
			&opening.JobTitle,
			&opening.CompanyDomain,
			&opening.CompanyName,
			&cities,
			&opening.PaginationKey,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning opening row: %w", err)
		}
		if len(cities) > 0 {
			opening.Cities = cities
		}
		openings = append(openings, opening)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating opening rows: %w", err)
	}

	return openings, nil
}
