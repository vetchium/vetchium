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
SELECT a.id, a.cover_letter, a.created_at, 
       h.handle as hub_user_handle,
       (SELECT domain_name FROM hub_users_official_emails hue
        JOIN domains d ON d.domain_name = split_part(hue.official_email, '@', 2)
        WHERE hue.hub_user_id = a.hub_user_id
        ORDER BY hue.created_at DESC
        LIMIT 1) as last_employer_domain,
       a.application_state, a.color_tag
FROM applications a
JOIN hub_users h ON h.id = a.hub_user_id
WHERE a.id = $1 AND a.employer_id = $2
`

	var app employer.Application
	var coverLetter *string
	var lastEmployerDomain *string
	var colorTag *string

	err := p.pool.QueryRow(ctx, query, applicationID, orgUser.EmployerID).
		Scan(&app.ID, &coverLetter, &app.CreatedAt,
			&app.HubUserHandle, &lastEmployerDomain,
			&app.State, &colorTag)
	if err != nil {
		p.log.Err("failed to get application", "error", err)
		return employer.Application{}, db.ErrInternal
	}

	app.CoverLetter = coverLetter
	app.HubUserLastEmployerDomain = lastEmployerDomain
	if colorTag != nil {
		tag := employer.ApplicationColorTag(*colorTag)
		app.ColorTag = &tag
	}

	return app, nil
}
