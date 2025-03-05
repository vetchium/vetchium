package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) InviteHubUser(
	ctx context.Context,
	inviteHubUserReq db.InviteHubUserReq,
) error {
	return nil
}
