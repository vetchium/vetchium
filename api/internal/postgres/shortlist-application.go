package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) ShortlistApplication(
	ctx context.Context,
	shortlistRequest db.ShortlistRequest,
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
	defer tx.Rollback(context.Background())

	candidacyQuery := `
INSERT INTO candidacies (id, application_id, employer_id, opening_id, created_by, candidacy_state)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
	`
	var candidacyID string
	err = tx.QueryRow(
		ctx,
		candidacyQuery,
		shortlistRequest.CandidacyID,
		shortlistRequest.ApplicationID,
		orgUser.EmployerID,
		shortlistRequest.OpeningID,
		orgUser.ID,
		common.InterviewingCandidacyState,
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
            WHERE id = $1 AND employer_id = $2
        ) THEN $5
        WHEN EXISTS (
            SELECT 1 FROM applications 
            WHERE id = $1 AND employer_id = $2 
            AND application_state != $4
        ) THEN $6
        ELSE $7
    END as status
),
update_result AS (
    UPDATE applications
    SET application_state = $3
    WHERE id = $1 
    AND employer_id = $2 
    AND application_state = $4
    AND (SELECT status FROM application_check) = $7
)
SELECT status FROM application_check;
`
	var status string
	err = tx.QueryRow(
		ctx,
		applicationQuery,
		shortlistRequest.ApplicationID,
		orgUser.EmployerID,
		common.ShortlistedAppState,
		common.AppliedAppState,
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
		p.log.Dbg("application not found", "id", shortlistRequest.ApplicationID)
		return db.ErrNoApplication
	case statusWrongState:
		p.log.Dbg(
			"application is in wrong state",
			"id",
			shortlistRequest.ApplicationID,
		)
		return db.ErrApplicationStateInCompatible
	case statusOK:
		// Continue with the rest of the function
	default:
		p.log.Err("unexpected status", "status", status)
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
		ctx,
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

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
