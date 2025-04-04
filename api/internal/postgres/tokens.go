package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) CreateOrgUserToken(
	ctx context.Context,
	tokenReq db.EmployerTokenReq,
) error {
	query := `
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`
	_, err := p.pool.Exec(
		ctx,
		query,
		tokenReq.Token,
		tokenReq.OrgUserID,
		tokenReq.ValidityDuration.Minutes(),
		tokenReq.TokenType,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "org_user_tokens_pkey" {
			p.log.Err("duplicate token generated", "error", err)
			return err
		}
		p.log.Err("failed to create org user token", "error", err)
		return err
	}

	return nil
}
