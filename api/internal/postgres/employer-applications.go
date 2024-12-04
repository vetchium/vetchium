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

func (p *PG) SetApplicationColorTag(
	c context.Context,
	req vetchi.SetApplicationColorTagRequest,
) error {
	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
	)

	orgUser, ok := c.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
		WITH application_check AS (
			SELECT CASE
				WHEN NOT EXISTS (
					SELECT 1 FROM applications
					WHERE id = $2 AND employer_id = $3
				) THEN $5
				WHEN EXISTS (
					SELECT 1 FROM applications
					WHERE id = $2 AND employer_id = $3
					AND application_state != $4
				) THEN $6
				ELSE $7
			END as status
		)
		UPDATE applications
		SET color_tag = $1
		WHERE id = $2
		AND employer_id = $3
		AND application_state = $4
		AND (SELECT status FROM application_check) = $7
		RETURNING (SELECT status FROM application_check);
	`

	var status string
	err := p.pool.QueryRow(
		c,
		query,
		req.ColorTag,
		req.ApplicationID,
		orgUser.EmployerID,
		vetchi.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to add application color tag", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		return nil
	default:
		p.log.Err("failed to add application color tag", "error", err)
		return db.ErrInternal
	}
}

func (p *PG) RemoveApplicationColorTag(
	c context.Context,
	req vetchi.RemoveApplicationColorTagRequest,
) error {
	orgUser, ok := c.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
	)

	query := `
WITH application_check AS (
	SELECT CASE
		WHEN NOT EXISTS (
			SELECT 1 FROM applications 
			WHERE id = $1 AND employer_id = $2
		) THEN $4
		WHEN EXISTS (
			SELECT 1 FROM applications
			WHERE id = $1 AND employer_id = $2
			AND application_state != $3
		) THEN $5
		ELSE $6
	END as status
)
UPDATE applications
SET color_tag = NULL
WHERE id = $1
AND employer_id = $2
AND application_state = $3
AND (SELECT status FROM application_check) = $6
RETURNING (SELECT status FROM application_check);
`

	var status string
	err := p.pool.QueryRow(
		c,
		query,
		req.ApplicationID,
		orgUser.EmployerID,
		vetchi.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to remove application color tag", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		return nil
	default:
		p.log.Err("unexpected status when removing color tag", "error", err)
		return db.ErrInternal
	}
}

func (p *PG) ShortlistApplication(
	c context.Context,
	req db.ShortlistRequest,
) error {
	return nil
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
			d.domain_name as primary_domain
		FROM applications a
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
