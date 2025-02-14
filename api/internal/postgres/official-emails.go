package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		// Check for unique constraint violation on official_email
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			p.log.Dbg("duplicate official email", "email", req.Email.EmailTo[0])
			return db.ErrDuplicateOfficialEmail
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

func (p *PG) GetOfficialEmail(
	ctx context.Context,
	email string,
) (*db.OfficialEmail, error) {
	var result db.OfficialEmail
	var lastVerifiedAt sql.NullTime

	err := p.pool.QueryRow(ctx, `
		SELECT
			official_email,
			last_verified_at,
			verification_code IS NOT NULL as verify_in_progress
		FROM hub_users_official_emails
		WHERE official_email = $1
	`, email).Scan(&result.Email, &lastVerifiedAt, &result.VerifyInProgress)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, db.ErrOfficialEmailNotFound
		}
		return nil, fmt.Errorf("failed to get official email: %w", err)
	}

	if lastVerifiedAt.Valid {
		result.LastVerifiedAt = &lastVerifiedAt.Time
	}

	return &result, nil
}

func (p *PG) UpdateOfficialEmailVerificationCode(
	ctx context.Context,
	req db.UpdateOfficialEmailVerificationCodeReq,
) error {
	commandTag, err := p.pool.Exec(ctx, `
		UPDATE hub_users_official_emails
		SET
			verification_code = $1,
			verification_code_expires_at = timezone('UTC', now()) + interval '24 hours'
		WHERE official_email = $2 AND hub_user_id = $3
	`, req.Code, req.Email, req.HubUser.ID)

	if err != nil {
		return fmt.Errorf("failed to update verification code: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return db.ErrOfficialEmailNotFound
	}

	return nil
}

func (p *PG) VerifyOfficialEmail(
	ctx context.Context,
	email string,
	code string,
) error {
	// Update the verification status and clear the verification code
	commandTag, err := p.pool.Exec(ctx, `
		UPDATE hub_users_official_emails
		SET
			verification_code = NULL,
			verification_code_expires_at = NULL,
			last_verified_at = timezone('UTC', now())
		WHERE official_email = $1
		  AND verification_code = $2
		  AND verification_code_expires_at > timezone('UTC', now())
	`, email, code)

	if err != nil {
		return fmt.Errorf("failed to verify official email: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		// Check if the email exists
		var exists bool
		err := p.pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM hub_users_official_emails
				WHERE official_email = $1
			)
		`, email).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}

		if !exists {
			return db.ErrOfficialEmailNotFound
		}
		return db.ErrInvalidVerificationCode
	}

	return nil
}

func (p *PG) PruneOfficialEmailCodes(ctx context.Context) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE hub_users_official_emails
		SET
			verification_code = NULL,
			verification_code_expires_at = NULL
		WHERE verification_code_expires_at < timezone('UTC', now())
		AND verification_code IS NOT NULL
	`)
	if err != nil {
		return fmt.Errorf("failed to prune official email codes: %w", err)
	}

	return nil
}

func (p *PG) DeleteOfficialEmail(ctx context.Context, email string) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	commandTag, err := p.pool.Exec(ctx, `
		DELETE FROM hub_users_official_emails
		WHERE official_email = $1 AND hub_user_id = $2
	`, email, hubUserID)

	if err != nil {
		return fmt.Errorf("failed to delete official email: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return db.ErrOfficialEmailNotFound
	}

	return nil
}
