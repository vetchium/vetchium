package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetResumeDetails(
	ctx context.Context,
	request employer.GetResumeRequest,
) (db.ResumeDetails, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ResumeDetails{}, db.ErrInternal
	}

	query := `
SELECT a.resume_sha, h.handle, a.id
FROM applications a
JOIN hub_users h ON h.id = a.hub_user_id
WHERE a.id = $1 AND a.employer_id = $2
`
	var details db.ResumeDetails
	err := p.pool.QueryRow(ctx, query, request.ApplicationID, orgUser.EmployerID).
		Scan(&details.SHA, &details.HubUserHandle, &details.ApplicationID)
	if err != nil {
		p.log.Err("failed to get resume details", "error", err)
		return db.ResumeDetails{}, db.ErrInternal
	}

	return details, nil
}

func (p *PG) GetApplication(
	ctx context.Context,
	applicationID string,
) (employer.Application, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.Application{}, db.ErrInternal
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
SELECT a.id, a.cover_letter, a.created_at, 
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
WHERE a.id = $1 AND a.employer_id = $2
`

	var app employer.Application
	var coverLetter *string
	var colorTag *string

	err := p.pool.QueryRow(ctx, query, applicationID, orgUser.EmployerID).
		Scan(&app.ID, &coverLetter, &app.CreatedAt,
			&app.HubUserHandle, &app.HubUserName, &app.HubUserShortBio,
			&app.HubUserLastEmployerDomains, &app.State, &colorTag,
			&app.Endorsers)
	if err != nil {
		p.log.Err("failed to get application", "error", err)
		return employer.Application{}, db.ErrInternal
	}

	app.CoverLetter = coverLetter
	if colorTag != nil {
		tag := employer.ApplicationColorTag(*colorTag)
		app.ColorTag = &tag
	}

	return app, nil
}
