package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/employer"
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

	// TODO: We need a better way to ORDER BY the applications
	// based on the scores across multiple models
	query := `
		WITH endorsed_applications AS (
			SELECT
				a.id,
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
			FROM applications a
			JOIN application_endorsements ae ON a.id = ae.application_id
			JOIN hub_users h ON ae.endorser_id = h.id
			WHERE ae.state = 'ENDORSED'
			GROUP BY a.id
		),
		application_model_scores AS (
			SELECT
				application_id,
				jsonb_agg(
					jsonb_build_object(
						'model_name', model_name,
						'score', score
					)
				) as scores
			FROM application_scores
			GROUP BY application_id
		)
		SELECT
			a.id,
			a.cover_letter,
			a.created_at,
			h.handle as hub_user_handle,
			h.full_name as hub_user_name,
			h.short_bio as hub_user_short_bio,
			(
				SELECT array_agg(d.domain_name ORDER BY hue.last_verified_at DESC)
				FROM (
					SELECT DISTINCT ON (hub_user_id) hub_user_id, domain_id, last_verified_at
					FROM hub_users_official_emails
					WHERE hub_user_id = h.id
					AND last_verified_at IS NOT NULL
					ORDER BY hub_user_id, last_verified_at DESC
				) hue
				JOIN domains d ON d.id = hue.domain_id
			) as hub_user_last_employer_domains,
			a.application_state,
			a.color_tag,
			COALESCE(ea.endorsers, '[]'::jsonb) as endorsers,
			COALESCE(ams.scores, '[]'::jsonb) as scores
		FROM applications a
		JOIN hub_users h ON h.id = a.hub_user_id
		LEFT JOIN endorsed_applications ea ON ea.id = a.id
		LEFT JOIN application_model_scores ams ON ams.application_id = a.id
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
			&app.Scores,
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
