package postgres

import (
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) AddOfficialEmail(req db.AddOfficialEmailReq) error {
	ctx := req.Context

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Add the official email to the hub_users_official_emails table
	officialEmailsQuery := `
		WITH domain_lookup AS (
			SELECT id FROM domains WHERE domain_name = $2
		)
		INSERT INTO hub_users_official_emails (
			hub_user_id,
			domain_id,
			official_email,
			verification_code,
			verification_code_expires_at
		)
		SELECT
			$1,
			domain_lookup.id,
			$3,
			$4,
			timezone('UTC', now()) + ($5 * INTERVAL '1 minute')
		FROM domain_lookup
		WHERE EXISTS (SELECT 1 FROM domain_lookup)
		RETURNING verification_code`

	verificationCodeExpiry := 24 * time.Hour
	domain := req.Email.EmailTo[0][strings.Index(req.Email.EmailTo[0], "@")+1:]
	var verificationCode string
	err = tx.QueryRow(
		ctx,
		officialEmailsQuery,
		req.HubUser.ID,
		domain,
		req.Email.EmailTo[0],
		req.Code,
		verificationCodeExpiry.Minutes(),
	).Scan(&verificationCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("domain not found", "domain", domain)
			return db.ErrNoEmployer // Domain not found implies employer not found
		}
		p.log.Err("failed to add official email", "error", err)
		return err
	}

	// Send the email with the token to the added official email address
	tokenMailQuery := `
		INSERT INTO emails (
			email_from,
			email_to,
			email_subject,
			email_html_body,
			email_text_body,
			email_state
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		RETURNING email_key`

	var emailKey string
	err = tx.QueryRow(ctx, tokenMailQuery,
		req.Email.EmailFrom,
		req.Email.EmailTo,
		req.Email.EmailSubject,
		req.Email.EmailHTMLBody,
		req.Email.EmailTextBody,
		"PENDING",
	).Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to send token mail", "error", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
