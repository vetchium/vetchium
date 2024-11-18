package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (p *PG) GetHubUserByToken(
	ctx context.Context,
	tfaToken string,
	tfaCode string,
) (db.HubUserTO, error) {
	query := `
SELECT
    *
FROM
    hub_users hub
    JOIN hub_user_tokens tfa_token ON tfa_token.hub_user_id = hub.id
        AND tfa_token.token = $1
        AND tfa_token.token_type = $2
    JOIN hub_user_tokens tfa_code ON tfa_code.hub_user_id = hub.id
        AND tfa_code.token = $3
        AND tfa_code.token_type = $4
`

	var hubUser db.HubUserTO
	if err := p.pool.QueryRow(ctx, query, tfaToken, db.HubUserTFAToken, tfaCode, db.HubUserTFACode).Scan(
		&hubUser.ID,
		&hubUser.FullName,
		&hubUser.Handle,
		&hubUser.Email,
		&hubUser.PasswordHash,
		&hubUser.CreatedAt,
		&hubUser.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.HubUserTO{}, db.ErrNoHubUser
		}

		p.log.Err("failed to get hub user", "error", err)
		return db.HubUserTO{}, db.ErrInternal
	}

	return hubUser, nil
}

func (p *PG) CreateHubUserToken(
	ctx context.Context,
	tokenReq db.HubTokenReq,
) error {
	query := `
INSERT INTO hub_user_tokens(token, hub_user_id, token_valid_till, token_type) VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`
	_, err := p.pool.Exec(
		ctx,
		query,
		tokenReq.Token,
		tokenReq.HubUserID,
		tokenReq.ValidityDuration.Minutes(),
		tokenReq.TokenType,
	)
	if err != nil {
		p.log.Err("failed to create hub user token", "error", err)
		return err
	}

	return nil
}
