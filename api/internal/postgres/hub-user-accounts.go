package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) SignupHubUser(ctx context.Context, req db.SignupHubUserReq) error {
	p.log.Dbg("Entered SignupHubUser", "request", req)

	const (
		statusUserExists             = "USER_EXISTS"
		statusInviteExists           = "INVITE_EXISTS"
		statusDomainNotApproved      = "DOMAIN_NOT_APPROVED"
		statusCanSignup              = "CAN_SIGNUP"
		statusSignupInviteCreated    = "SIGNUP_INVITE_CREATED"
		statusSignupInviteNotCreated = "SIGNUP_INVITE_NOT_CREATED"
	)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(context.Background())

	// Combined check and insert query
	inviteQuery := `
WITH eligibility_check AS (
    SELECT check_hub_user_signup_eligibility($1) as status
),
insert_invite AS (
    INSERT INTO hub_user_invites (email, token, token_valid_till)
    SELECT $1, $2, $3
    WHERE (SELECT status FROM eligibility_check) = '` + statusCanSignup + `'
    RETURNING email
)
SELECT COALESCE(
    (SELECT '` + statusSignupInviteCreated + `' FROM insert_invite),
    (SELECT status FROM eligibility_check)
) as final_status;
`
	var finalStatus string
	err = tx.QueryRow(
		ctx,
		inviteQuery,
		req.EmailAddress,
		req.Token,
		req.TokenValidTill,
	).Scan(&finalStatus)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			// This case should ideally be caught by check_hub_user_signup_eligibility
			// but as a fallback or if unique_token constraint is hit.
			if pgErr.ConstraintName == "unique_user" {
				p.log.Dbg(
					"invite already exists (caught by DB constraint)",
					"email",
					req.EmailAddress,
				)
				return db.ErrInviteNotNeeded
			} else if pgErr.ConstraintName == "unique_token" {
				p.log.Err("duplicate token generated for signup invite", "token", req.Token, "error", err)
				return db.ErrInternal // Or a more specific error like ErrTokenGenerationFailed
			}
		}
		p.log.Err(
			"failed during signup eligibility check or invite insert",
			"email",
			req.EmailAddress,
			"error",
			err,
		)
		return db.ErrInternal
	}

	switch finalStatus {
	case statusSignupInviteCreated:
		p.log.Dbg("signup invite created", "email", req.EmailAddress)
		// Proceed to insert email
	case statusUserExists:
		p.log.Dbg("user already exists", "email", req.EmailAddress)
		return db.ErrInviteNotNeeded
	case statusInviteExists:
		p.log.Dbg("invite already exists", "email", req.EmailAddress)
		return db.ErrInviteNotNeeded
	case statusDomainNotApproved:
		extractedDomain := "unknown"
		emailParts := strings.Split(req.EmailAddress, "@")
		if len(emailParts) == 2 {
			extractedDomain = emailParts[1]
		}
		p.log.Dbg("domain not approved for signup",
			"email", req.EmailAddress,
			"domain", extractedDomain,
		)
		return db.ErrDomainNotApprovedForSignup
	default:
		p.log.Err("unexpected status from signup eligibility check",
			"status", finalStatus,
			"email", req.EmailAddress,
		)
		return db.ErrInternal
	}

	// Insert into emails table if invite was created
	emailQuery := `
INSERT INTO emails (email_from, email_to, email_cc, email_subject, email_html_body, email_text_body, email_state)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		req.InviteMail.EmailFrom,
		req.InviteMail.EmailTo,
		req.InviteMail.EmailCC,
		req.InviteMail.EmailSubject,
		req.InviteMail.EmailHTMLBody,
		req.InviteMail.EmailTextBody,
		req.InviteMail.EmailState,
	)
	if err != nil {
		p.log.Err("failed to create email for signup invite", "error", err)
		// Note: The invite is already in the hub_user_invites table.
		// Depending on desired behavior, we might want to attempt to rollback the invite insert
		// or mark the invite as needing email resend. For now, we return internal error.
		return db.ErrInternal
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction for signup", "error", err)
		return db.ErrInternal
	}

	p.log.Dbg("SignupHubUser completed successfully", "email", req.EmailAddress)
	return nil
}

func (p *PG) ChangeEmailAddress(ctx context.Context, email string) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	_, err = p.pool.Exec(
		ctx,
		"UPDATE hub_users SET email = $1 WHERE id = $2",
		email,
		hubUserID,
	)
	if err != nil {
		p.log.Err("failed to update email", "error", err)
		return db.ErrInternal
	}

	p.log.Dbg("email updated successfully", "email", email)
	return nil
}
