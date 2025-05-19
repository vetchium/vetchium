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
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) AddOfficialEmail(req db.AddOfficialEmailReq) error {
	ctx := req.Context

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction for AddOfficialEmail",
			"error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(context.Background())

	// Check if user has reached the maximum allowed official emails
	countQuery := `SELECT COUNT(*) FROM hub_users_official_emails WHERE hub_user_id = $1`
	var emailCount int
	err = tx.QueryRow(ctx, countQuery, req.HubUser.ID).Scan(&emailCount)
	if err != nil {
		p.log.Err("failed to count official emails", "error", err,
			"hub_user_id", req.HubUser.ID)
		return db.ErrInternal
	}

	if emailCount >= 50 {
		p.log.Dbg("user has reached max official emails", "count", emailCount,
			"hub_user_id", req.HubUser.ID)
		return db.ErrTooManyOfficialEmails
	}

	// Extract domain and ensure it's valid before proceeding
	domainName := extractDomainFromEmail(req.Email.EmailTo[0])
	if domainName == "" {
		p.log.Err("could not extract domain from email",
			"email", req.Email.EmailTo[0])
		return db.ErrInternal
	}

	// Get domain_id. If domain does not exist, create dummy employer and domain.
	var domainID string
	queryDomainID := `SELECT id FROM domains WHERE domain_name = $1`
	err = tx.QueryRow(ctx, queryDomainID, domainName).Scan(&domainID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("Domain not found, creating dummy", "domain", domainName)
			var dummyEmployerID string
			createDummySQL := `SELECT get_or_create_dummy_employer($1)`
			errCreate := tx.QueryRow(ctx, createDummySQL, domainName).
				Scan(&dummyEmployerID)
			if errCreate != nil {
				p.log.Err("Failed to exec get_or_create_dummy_employer",
					"error", errCreate, "domain", domainName)
				return db.ErrInternal
			}
			// Fetch the newly created domain's ID
			errFetchDomain := tx.QueryRow(ctx, queryDomainID, domainName).
				Scan(&domainID)
			if errFetchDomain != nil {
				p.log.Err("Failed to fetch domain_id after dummy creation",
					"error", errFetchDomain, "domain", domainName)
				return db.ErrInternal
			}
			p.log.Dbg(
				"Created and fetched new domain_id",
				"domain_id",
				domainID,
				"employer_id",
				dummyEmployerID,
			)
		} else {
			p.log.Err("failed to query domain_id", "error", err, "domain", domainName)
			return db.ErrInternal
		}
	}

	officialEmailsQuery := `
		INSERT INTO hub_users_official_emails (
			hub_user_id,
			domain_id,
			official_email,
			verification_code,
			verification_code_expires_at
		)
		VALUES ($1, $2, $3, $4, timezone('UTC', now()) + ($5 * INTERVAL '1 minute'))
		RETURNING verification_code`

	verificationCodeExpiry := 24 * time.Hour
	var verificationCode string

	err = tx.QueryRow(
		ctx,
		officialEmailsQuery,
		req.HubUser.ID,
		domainID,
		req.Email.EmailTo[0],
		req.Code,
		verificationCodeExpiry.Minutes(),
	).Scan(&verificationCode)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "hub_users_official_emails_pkey" ||
				strings.Contains(pgErr.Message, "official_email") {
				p.log.Dbg(
					"duplicate official email",
					"email",
					req.Email.EmailTo[0],
					"hub_user_id",
					req.HubUser.ID,
				)
				return db.ErrDuplicateOfficialEmail
			}
		}
		p.log.Err(
			"failed to insert official email record",
			"error",
			err,
			"hub_user_id",
			req.HubUser.ID,
			"domain_id",
			domainID,
			"email",
			req.Email.EmailTo[0],
		)
		return db.ErrInternal
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
	p.log.Dbg("Enqueued token mail", "email_key", emailKey)

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	p.log.Dbg("Added official email", "email", req.Email.EmailTo)
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
		p.log.Err("failed to verify official email", "error", err)
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
			p.log.Err("failed to check email existence", "error", err)
			return fmt.Errorf("failed to check email existence: %w", err)
		}

		if !exists {
			p.log.Dbg("email does not exist", "email", email)
			return db.ErrOfficialEmailNotFound
		}
		p.log.Dbg("invalid verification code", "email", email)
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
