package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) GetHubUserByEmail(
	ctx context.Context,
	email string,
) (db.HubUserTO, error) {
	query := "SELECT * FROM hub_users WHERE email = $1"

	var hubUser db.HubUserTO
	if err := p.pool.QueryRow(ctx, query, email).Scan(
		&hubUser.ID,
		&hubUser.FullName,
		&hubUser.Handle,
		&hubUser.Email,
		&hubUser.PasswordHash,
		&hubUser.CreatedAt,
		&hubUser.UpdatedAt,
	); err != nil {
		return db.HubUserTO{}, err
	}

	return hubUser, nil
}

func (p *PG) InitHubUserTFA(
	ctx context.Context,
	tfa db.HubUserTFA,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tfaTokenQuery := `
INSERT INTO hub_user_tokens(token, hub_user_id, token_valid_till, token_type) VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)`
	_, err = tx.Exec(
		ctx,
		tfaTokenQuery,
		tfa.TFAToken.Token,
		tfa.TFAToken.HubUserID,
		tfa.TFAToken.ValidityDuration.Minutes(),
		db.HubUserTFAToken,
	)
	if err != nil {
		p.log.Err("failed to insert TGT", "error", err)
		return err
	}

	tfaCodeQuery := `
INSERT INTO hub_user_tokens(token, hub_user_id, token_valid_till, token_type) VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)`
	_, err = tx.Exec(
		ctx,
		tfaCodeQuery,
		tfa.TFACode.Token,
		tfa.TFACode.HubUserID,
		tfa.TFACode.ValidityDuration.Minutes(),
		db.HubUserTFACode,
	)
	if err != nil {
		p.log.Err("failed to insert TFA code", "error", err)
		return err
	}

	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		tfa.Email.EmailFrom,
		tfa.Email.EmailTo,
		tfa.Email.EmailSubject,
		tfa.Email.EmailHTMLBody,
		tfa.Email.EmailTextBody,
		tfa.Email.EmailState,
	)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
