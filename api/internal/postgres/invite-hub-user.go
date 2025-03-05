package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) InviteHubUser(
	ctx context.Context,
	inviteHubUserReq db.InviteHubUserReq,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to start transaction", "err", err)
		return db.ErrInternal
	}

	defer tx.Rollback(context.Background())

	// TODO: Insert the inviteHubUserReq.InviteEmail into the emails table and
	// create a new record on the hub_user_invites table with the email_id
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

	// TODO: Fix SQL query to ensure that the user does not exist already
	// in hub_users or an invite already sent
	inviteHubUserQuery := `
INSERT INTO hub_user_invites (email_id, token, valid_till)
	VALUES ($1, $2, NOW() + INTERVAL '3 days')`
	_, err = tx.Exec(
		ctx,
		inviteHubUserQuery,
		string(inviteHubUserReq.EmailAddress),
		inviteHubUserReq.Token,
	)
	if err != nil {
		p.log.Err("failed to create hub user invite", "error", err)
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
