package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) SignupHubUser(ctx context.Context, req db.SignupHubUserReq) error {
	p.log.Dbg("Entered SignupHubUser", "request", req)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	inviteQuery := `
INSERT INTO hub_user_invites (email, token, token_valid_till)
VALUES ($1, $2, $3)
`
	_, err = tx.Exec(
		ctx,
		inviteQuery,
		req.EmailAddress,
		req.Token,
		req.TokenValidTill,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "unique_user" {
			p.log.Dbg("invite already exists", "email", req.EmailAddress)
			return db.ErrInviteNotNeeded
		}
		p.log.Err("failed to insert invite", "error", err)
		return err
	}

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
		p.log.Err("failed to create email", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

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
