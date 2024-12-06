package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) GetApplicationsForEmployer(
	c context.Context,
	req vetchi.GetApplicationsRequest,
) ([]vetchi.Application, error) {
	orgUser, ok := c.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	query := `
		SELECT
			a.id,
			a.cover_letter,
			a.created_at,
			a.original_filename,
			a.internal_filename,
			h.handle as hub_user_handle,
			a.application_state
		FROM applications a
		JOIN hub_users h ON h.id = a.hub_user_id
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
		query += ` AND (h.handle ILIKE $4 OR h.full_name ILIKE $4)`
		args = append(args, "%"+*req.SearchQuery+"%")
	}

	// Add pagination if key is provided
	if req.PaginationKey != nil {
		query += ` AND a.id > $5`
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

	var applications []vetchi.Application
	for rows.Next() {
		var app vetchi.Application
		var internalFilename string

		err := rows.Scan(
			&app.ID,
			&app.CoverLetter,
			&app.CreatedAt,
			&app.Filename,
			&internalFilename,
			&app.HubUserHandle,
			&app.State,
		)
		if err != nil {
			p.log.Err("failed to scan application", "error", err)
			return nil, db.ErrInternal
		}

		// Set the resume URL using the internal filename
		app.Resume = "/resumes/" + internalFilename

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
