package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) RejectApplication(
	ctx context.Context,
	rejectRequest db.RejectApplicationRequest,
) error {
	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
	)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

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
RETURNING (SELECT status FROM application_check)
`

	var status string
	err = tx.QueryRow(
		ctx,
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
		ctx,
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

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
