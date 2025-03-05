package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) InviteHubUser(
	ctx context.Context,
	inviteHubUserReq db.InviteHubUserReq,
) error {
	const (
		statusUserExists   = "USER_EXISTS"
		statusInviteExists = "INVITE_EXISTS"
		statusCanInsert    = "CAN_INSERT"
		statusInserted     = "INSERTED"
	)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to start transaction", "err", err)
		return db.ErrInternal
	}

	defer tx.Rollback(context.Background())

	emailQuery := `
INSERT INTO emails (email_from, email_to, email_cc, email_subject, email_html_body, email_text_body, email_state)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		inviteHubUserReq.InviteMail.EmailFrom,
		inviteHubUserReq.InviteMail.EmailTo,
		inviteHubUserReq.InviteMail.EmailCC,
		inviteHubUserReq.InviteMail.EmailSubject,
		inviteHubUserReq.InviteMail.EmailHTMLBody,
		inviteHubUserReq.InviteMail.EmailTextBody,
		inviteHubUserReq.InviteMail.EmailState,
	)
	if err != nil {
		p.log.Err("failed to create email", "error", err)
		return db.ErrInternal
	}

	inviteHubUserQuery := `
WITH checks AS (
    SELECT
        EXISTS(SELECT 1 FROM hub_users WHERE email = $1) as user_exists,
        EXISTS(SELECT 1 FROM hub_user_invites WHERE email = $1) as invite_exists
),
status AS (
    SELECT CASE
        WHEN user_exists THEN '` + statusUserExists + `'
        WHEN invite_exists THEN '` + statusInviteExists + `'
        ELSE '` + statusCanInsert + `'
    END as check_status
    FROM checks
),
ins AS (
    INSERT INTO hub_user_invites (email, token, token_valid_till)
    SELECT $1, $2, NOW() + INTERVAL '3 days'
    WHERE EXISTS (SELECT 1 FROM status WHERE check_status = '` + statusCanInsert + `')
    RETURNING email
)
SELECT COALESCE(
    (SELECT '` + statusInserted + `' FROM ins),
    (SELECT check_status FROM status)
);
`
	var status string
	err = tx.QueryRow(
		ctx,
		inviteHubUserQuery,
		string(inviteHubUserReq.EmailAddress),
		inviteHubUserReq.Token,
	).Scan(&status)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "unique_token":
				p.log.Err("duplicate token generated", "error", err)
				return db.ErrInternal
			default:
				p.log.Err("unexpected constraint violation", "error", err)
				return db.ErrInternal
			}
		}
		p.log.Err("failed to create hub user invite", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusInserted:
		// Success case
	case statusUserExists:
		return db.ErrInviteNotNeeded
	case statusInviteExists:
		return db.ErrInviteNotNeeded
	default:
		p.log.Err("unexpected status", "status", status)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "err", err)
		return db.ErrInternal
	}

	p.log.Dbg("hub user invited", "inviteHubUserReq", inviteHubUserReq)
	return nil
}
