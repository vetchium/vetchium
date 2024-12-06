package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) ShortlistApplication(
	ctx context.Context,
	shortlistRequest db.ShortlistRequest,
) error {
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
		vetchi.InterviewingCandidacyState,
	).
		Scan(&candidacyID)
	if err != nil {
		p.log.Err("failed to insert candidacy", "error", err)
		return db.ErrInternal
	}

	applicationQuery := `
UPDATE applications
SET application_state = $1
WHERE id = $2 AND employer_id = $3 AND application_state = $4
`
	result, err := tx.Exec(
		ctx,
		applicationQuery,
		vetchi.ShortlistedAppState,
		shortlistRequest.ApplicationID,
		orgUser.EmployerID,
		vetchi.AppliedAppState,
	)
	if err != nil {
		p.log.Err("failed to update application", "error", err)
		return db.ErrInternal
	}

	if result.RowsAffected() == 0 {
		p.log.Err("application not found", "id", shortlistRequest.ApplicationID)
		return db.ErrNoApplication
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

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
