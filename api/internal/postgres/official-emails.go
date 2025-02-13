package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) AddOfficialEmail(req db.AddOfficialEmailReq) error {
	ctx := req.Context

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Check if user has reached the maximum allowed official emails
	countQuery := `
		SELECT COUNT(*) 
		FROM hub_users_official_emails 
		WHERE hub_user_id = $1`

	var emailCount int
	err = tx.QueryRow(ctx, countQuery, req.HubUser.ID).Scan(&emailCount)
	if err != nil {
		p.log.Err("failed to count official emails", "error", err)
		return err
	}

	if emailCount >= 50 {
		p.log.Dbg(
			"user has reached maximum allowed official emails",
			"count",
			emailCount,
		)
		return db.ErrTooManyOfficialEmails
	}

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
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("domain not found", "domain", domain)
			return db.ErrNoEmployer
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
		req.Email.EmailState,
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

func (p *PG) GetMyOfficialEmails(
	ctx context.Context,
) ([]hub.OfficialEmail, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return nil, err
	}

	query := `
		SELECT
			he.official_email,
			he.last_verified_at,
			he.verification_code IS NOT NULL AS verify_in_progress
		FROM hub_users_official_emails he
		WHERE he.hub_user_id = $1
	`

	rows, err := p.pool.Query(ctx, query, hubUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no official emails found", "hub_user_id", hubUserID)
			return []hub.OfficialEmail{}, nil
		}
		p.log.Err("failed to get my official emails", "error", err)
		return nil, err
	}
	defer rows.Close()

	emails := []hub.OfficialEmail{}
	for rows.Next() {
		var email hub.OfficialEmail
		if err := rows.Scan(&email.Email, &email.LastVerifiedAt, &email.VerifyInProgress); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over my official emails", "error", err)
		return nil, err
	}

	p.log.Dbg("my official emails", "emails", emails)
	return emails, nil
}
