package postgres

import (
	"context"
	"fmt"

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
