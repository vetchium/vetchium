package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetApplicationsForEmployer(
	c context.Context,
	req employer.GetApplicationsRequest,
) ([]employer.Application, error) {
	orgUser, ok := c.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	query := `
		WITH endorsed_applications AS (
			SELECT 
				ae.application_id,
				jsonb_agg(
					jsonb_build_object(
						'full_name', h.full_name,
						'short_bio', h.short_bio,
						'handle', h.handle,
						'current_company_domains', (
							SELECT array_agg(d.domain_name)
							FROM hub_users_official_emails hue
							JOIN domains d ON d.id = hue.domain_id
							WHERE hue.hub_user_id = h.id
							AND hue.last_verified_at IS NOT NULL
						)
					)
				) as endorsers
			FROM application_endorsements ae
			JOIN hub_users h ON h.id = ae.endorser_id
			WHERE ae.state = 'ENDORSED'
			GROUP BY ae.application_id
		)
		SELECT
			a.id,
			a.cover_letter,
			a.created_at,
			h.handle as hub_user_handle,
			h.full_name as hub_user_name,
			h.short_bio as hub_user_short_bio,
			(
				SELECT array_agg(d.domain_name)
				FROM hub_users_official_emails hue
				JOIN domains d ON d.id = hue.domain_id
				WHERE hue.hub_user_id = h.id
				AND hue.last_verified_at IS NOT NULL
				ORDER BY hue.last_verified_at DESC
				LIMIT 1
			) as hub_user_last_employer_domains,
			a.application_state,
			a.color_tag,
			COALESCE(ea.endorsers, '[]'::jsonb) as endorsers
		FROM applications a
		JOIN hub_users h ON h.id = a.hub_user_id
		LEFT JOIN endorsed_applications ea ON ea.application_id = a.id
		WHERE a.employer_id = $1
		AND a.opening_id = $2
		AND a.application_state = $3
	`

	args := []interface{}{
		orgUser.EmployerID,
		req.OpeningID,
		req.State,
	}

	if req.SearchQuery != nil {
		query += fmt.Sprintf(
			` AND (h.handle ILIKE $%d OR h.full_name ILIKE $%d)`,
			len(args)+1,
			len(args)+1,
		)
		args = append(args, "%"+*req.SearchQuery+"%")
	}

	if req.ColorTagFilter != nil {
		query += fmt.Sprintf(` AND a.color_tag = $%d`, len(args)+1)
		args = append(args, *req.ColorTagFilter)
	}

	// Add pagination if key is provided
	if req.PaginationKey != nil {
		query += fmt.Sprintf(` AND a.id > $%d`, len(args)+1)
		args = append(args, *req.PaginationKey)
	}

	// Add limit
	query += ` ORDER BY a.id LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Limit)

	rows, err := p.pool.Query(c, query, args...)
	if err != nil {
		p.log.Err("failed to query applications", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	applications := make([]employer.Application, 0)
	for rows.Next() {
		var app employer.Application

		err := rows.Scan(
			&app.ID,
			&app.CoverLetter,
			&app.CreatedAt,
			&app.HubUserHandle,
			&app.HubUserName,
			&app.HubUserShortBio,
			&app.HubUserLastEmployerDomains,
			&app.State,
			&app.ColorTag,
			&app.Endorsers,
		)
		if err != nil {
			p.log.Err("failed to scan application", "error", err)
			return nil, db.ErrInternal
		}

		applications = append(applications, app)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating applications", "error", err)
		return nil, db.ErrInternal
	}

	return applications, nil
}

func (p *PG) GetApplicationMailInfo(
	c context.Context,
	applicationID string,
) (db.ApplicationMailInfo, error) {
	orgUser, ok := c.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ApplicationMailInfo{}, db.ErrInternal
	}

	var mailInfo db.ApplicationMailInfo

	query := `
SELECT
	h.id as hub_user_id,
	h.state as hub_user_state,
	h.full_name,
	h.handle,
	h.email,
	h.preferred_language,
	e.id as employer_id,
	e.company_name,
	d.domain_name as primary_domain,
	o.id as opening_id,
	o.title as opening_title
FROM applications a
JOIN openings o ON a.opening_id = o.id
JOIN hub_users h ON h.id = a.hub_user_id
JOIN employers e ON e.id = a.employer_id
JOIN employer_primary_domains epd ON epd.employer_id = e.id
JOIN domains d ON d.id = epd.domain_id
WHERE a.id = $1
AND a.employer_id = $2
`

	err := p.pool.QueryRow(c, query, applicationID, orgUser.EmployerID).Scan(
		&mailInfo.HubUser.HubUserID,
		&mailInfo.HubUser.State,
		&mailInfo.HubUser.FullName,
		&mailInfo.HubUser.Handle,
		&mailInfo.HubUser.Email,
		&mailInfo.HubUser.PreferredLanguage,
		&mailInfo.Employer.EmployerID,
		&mailInfo.Employer.CompanyName,
		&mailInfo.Employer.PrimaryDomain,
		&mailInfo.Opening.OpeningID,
		&mailInfo.Opening.Title,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("either application not found or not owned by employer")
			return db.ApplicationMailInfo{}, db.ErrNoApplication
		}

		p.log.Err("failed to get application mail info", "error", err)
		return db.ApplicationMailInfo{}, db.ErrInternal
	}

	return mailInfo, nil
}
