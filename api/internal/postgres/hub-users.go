package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
)

func (p *PG) GetHubUserByEmail(
	ctx context.Context,
	email string,
) (db.HubUserTO, error) {
	query := `
SELECT
	hu.id,
	hu.full_name,
	hu.handle,
	hu.email,
	hu.state,
	hu.password_hash,
	hu.created_at,
	hu.updated_at
FROM hub_users hu WHERE email = $1`

	var hubUser db.HubUserTO
	err := p.pool.QueryRow(ctx, query, email).Scan(
		&hubUser.ID,
		&hubUser.FullName,
		&hubUser.Handle,
		&hubUser.Email,
		&hubUser.State,
		&hubUser.PasswordHash,
		&hubUser.CreatedAt,
		&hubUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.HubUserTO{}, db.ErrNoHubUser
		}
		p.log.Err("failed to get hub user", "error", err)
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
INSERT INTO
	hub_user_tokens(token, hub_user_id, token_valid_till, token_type)
VALUES
	($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`
	_, err = tx.Exec(
		ctx,
		tfaTokenQuery,
		tfa.TFAToken.Token,
		tfa.TFAToken.HubUserID,
		tfa.TFAToken.ValidityDuration.Minutes(),
		db.HubUserTFAToken,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "hub_user_tokens_pkey" {
			p.log.Err("duplicate token generated", "error", err)
			return err
		}
		p.log.Err("failed to insert TFA Token", "error", err)
		return err
	}

	tfaCodeQuery := `
INSERT INTO
	hub_user_tfa_codes(code, tfa_token)
VALUES
	($1, $2)
`
	_, err = tx.Exec(
		ctx,
		tfaCodeQuery,
		tfa.TFACode,
		tfa.TFAToken.Token,
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

func (p *PG) GetHubUserByTFACreds(
	ctx context.Context,
	tfaToken string,
	tfaCode string,
) (db.HubUserTO, error) {
	query := `
SELECT
	hu.id,
	hu.full_name,
	hu.handle,
	hu.email,
	hu.password_hash,
	hu.created_at,
	hu.updated_at
FROM
	hub_users hu,
	hub_user_tokens hut,
	hub_user_tfa_codes hutc
WHERE
	hutc.code = $1
	AND hut.token = $2
	AND hutc.tfa_token = hut.token
	AND hut.token_type = $3
	AND hu.id = hut.hub_user_id
`

	var hubUser db.HubUserTO
	err := p.pool.QueryRow(
		ctx,
		query,
		tfaCode,
		tfaToken,
		db.HubUserTFAToken,
	).Scan(
		&hubUser.ID,
		&hubUser.FullName,
		&hubUser.Handle,
		&hubUser.Email,
		&hubUser.PasswordHash,
		&hubUser.CreatedAt,
		&hubUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no hub user found", "error", err)
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

func (p *PG) Logout(ctx context.Context, token string) error {
	query := "DELETE FROM hub_user_tokens WHERE token = $1"

	_, err := p.pool.Exec(ctx, query, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// unlikely to happen
			p.log.Dbg("no hub user token found", "error", err)
			return nil
		}

		p.log.Err("failed to delete hub user token", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetMyHandle(ctx context.Context) (string, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return "", db.ErrInternal
	}

	return hubUser.Handle, nil
}

func (p *PG) AuthHubUser(
	ctx context.Context,
	token string,
) (db.HubUserTO, error) {
	query := `
SELECT
	hu.id,
	hu.full_name,
	hu.handle,
	hu.email,
	hu.password_hash,
	hu.created_at,
	hu.updated_at
FROM 
	hub_user_tokens hut,
	hub_users hu
WHERE
	hut.token = $1
	AND (hut.token_type = $2 OR hut.token_type = $3)
	AND hu.id = hut.hub_user_id
`

	var hubUser db.HubUserTO
	err := p.pool.QueryRow(ctx, query, token, db.HubUserSessionToken, db.HubUserLTSToken).
		Scan(
			&hubUser.ID,
			&hubUser.FullName,
			&hubUser.Handle,
			&hubUser.Email,
			&hubUser.PasswordHash,
			&hubUser.CreatedAt,
			&hubUser.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no hub user token found", "error", err)
			return db.HubUserTO{}, db.ErrNoHubUser
		}

		p.log.Err("failed to auth hub user", "error", err)
		return db.HubUserTO{}, db.ErrInternal
	}

	return hubUser, nil
}

func (p *PG) ChangeHubUserPassword(
	ctx context.Context,
	hubUserID uuid.UUID,
	newPasswordHash string,
) error {
	query := `
UPDATE hub_users 
SET password_hash = $1, updated_at = NOW() AT TIME ZONE 'utc'
WHERE id = $2`

	_, err := p.pool.Exec(ctx, query, newPasswordHash, hubUserID)
	if err != nil {
		p.log.Err("failed to update password", "error", err)
		return err
	}

	return nil
}

func (p *PG) InitHubUserPasswordReset(
	ctx context.Context,
	initPasswordResetReq db.HubUserInitPasswordReset,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tokensQuery := `
INSERT INTO hub_user_tokens(token, hub_user_id, token_valid_till, token_type) 
VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`
	_, err = tx.Exec(
		ctx,
		tokensQuery,
		initPasswordResetReq.Token,
		initPasswordResetReq.HubUserID,
		initPasswordResetReq.ValidityDuration.Minutes(),
		db.HubUserResetPasswordToken,
	)
	if err != nil {
		p.log.Err("failed to insert password reset token", "error", err)
		return err
	}

	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		initPasswordResetReq.Email.EmailFrom,
		initPasswordResetReq.Email.EmailTo,
		initPasswordResetReq.Email.EmailSubject,
		initPasswordResetReq.Email.EmailHTMLBody,
		initPasswordResetReq.Email.EmailTextBody,
		initPasswordResetReq.Email.EmailState,
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

func (p *PG) ResetHubUserPassword(
	ctx context.Context,
	hubUserPasswordReset db.HubUserPasswordReset,
) error {
	query := `
UPDATE
    hub_users
SET
    password_hash = $1,
    updated_at = NOW() AT TIME ZONE 'utc'
WHERE
    id = (
        SELECT
            hub_user_id
        FROM
            hub_user_tokens
        WHERE
            token = $2)
RETURNING id
`
	var hubUserID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		hubUserPasswordReset.PasswordHash,
		hubUserPasswordReset.Token,
	).Scan(&hubUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("invalid password reset token", "error", err)
			return db.ErrInvalidPasswordResetToken
		}

		p.log.Err("failed to reset password", "error", err)
		return err
	}

	p.log.Dbg("password reset", "hubUserID", hubUserID)

	return nil
}
