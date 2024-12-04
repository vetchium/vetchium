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

func (p *PG) RejectApplication(
	c context.Context,
	rejectRequest db.RejectApplicationRequest,
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

	tx, err := p.pool.Begin(c)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(c)

	// Modified query to check application existence and state
	applicationQuery := `
	WITH application_check AS (
		SELECT CASE
			WHEN NOT EXISTS (
				SELECT 1 FROM applications
				WHERE id = $2 AND employer_id = $3
			) THEN $4
			WHEN EXISTS (
				SELECT 1 FROM applications
				WHERE id = $2 AND employer_id = $3
				AND application_state != $1
			) THEN $5
			ELSE $6
		END as status
	)
	UPDATE applications
	SET application_state = $1
	WHERE id = $2 
	AND employer_id = $3 
	AND application_state = $7
	AND (SELECT status FROM application_check) = $6
	RETURNING (SELECT status FROM application_check);
	`

	var status string
	err = tx.QueryRow(
		c,
		applicationQuery,
		vetchi.RejectedAppState,
		rejectRequest.ApplicationID,
		orgUser.EmployerID,
		statusNotFound,
		statusWrongState,
		statusOK,
		vetchi.AppliedAppState,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to update application state", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		// continue with the rest of the function
	default:
		p.log.Err("unexpected status when updating application", "error", err)
		return db.ErrInternal
	}

	// Rest of the email insertion code remains the same
	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
	VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
	email_key
`
	var emailKey string
	err = tx.QueryRow(
		c,
		emailQuery,
		rejectRequest.Email.EmailFrom,
		rejectRequest.Email.EmailTo,
		rejectRequest.Email.EmailSubject,
		rejectRequest.Email.EmailHTMLBody,
		rejectRequest.Email.EmailTextBody,
		rejectRequest.Email.EmailState,
	).Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return db.ErrInternal
	}
	p.log.Dbg("RejectApplication email added", "email_key", emailKey)

	err = tx.Commit(c)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

func (p *PG) ShortlistApplication(
	c context.Context,
	shortlistRequest db.ShortlistRequest,
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

	tx, err := p.pool.Begin(c)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(c)

	candidacyQuery := `
INSERT INTO candidacies (id, application_id, employer_id, opening_id, created_by, candidacy_state)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
	`
	var candidacyID string
	err = tx.QueryRow(
		c,
		candidacyQuery,
		shortlistRequest.CandidacyID,
		shortlistRequest.ApplicationID,
		orgUser.EmployerID,
		shortlistRequest.OpeningID,
		orgUser.ID,
		vetchi.InterviewingCandidacyState,
	).
		Scan(&candidacyID)
	if err != nil {
		p.log.Err("failed to insert candidacy", "error", err)
		return db.ErrInternal
	}

	applicationQuery := `
	WITH application_check AS (
		SELECT CASE
			WHEN NOT EXISTS (
				SELECT 1 FROM applications
				WHERE id = $3 AND employer_id = $4
			) THEN $6
			WHEN EXISTS (
				SELECT 1 FROM applications
				WHERE id = $3 AND employer_id = $4
				AND application_state != $5
			) THEN $7
			ELSE $8
		END as status
	)
	UPDATE applications
	SET
		candidacy_id = $1,
		application_state = $2
	WHERE
		id = $3
		AND employer_id = $4
		AND application_state = $5
		AND (SELECT status FROM application_check) = $8
	RETURNING (SELECT status FROM application_check);
	`
	var status string
	err = tx.QueryRow(
		c,
		applicationQuery,
		candidacyID,
		vetchi.ShortlistedAppState,
		shortlistRequest.ApplicationID,
		orgUser.EmployerID,
		vetchi.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to update application", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		// continue with the rest of the function
	default:
		p.log.Err("unexpected status when updating application", "error", err)
		return db.ErrInternal
	}

	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    email_key
`
	var emailKey string
	err = tx.QueryRow(
		c,
		emailQuery,
		shortlistRequest.Email.EmailFrom,
		shortlistRequest.Email.EmailTo,
		shortlistRequest.Email.EmailSubject,
		shortlistRequest.Email.EmailHTMLBody,
		shortlistRequest.Email.EmailTextBody,
		shortlistRequest.Email.EmailState,
	).Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return db.ErrInternal
	}
	p.log.Dbg("ShortlistApplication email added", "email_key", emailKey)

	err = tx.Commit(c)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

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
