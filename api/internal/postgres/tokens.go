package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) CreateOrgUserToken(
	ctx context.Context,
	tokenReq db.TokenReq,
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
		// TODO: Check if the error is due to duplicate key value
		// and if so retry with a different token
		p.log.Error("failed to create org user token", "error", err)
		return err
	}

	return nil
}
