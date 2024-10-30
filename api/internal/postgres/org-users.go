package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) AddOrgUser(
	ctx context.Context,
	req db.AddOrgUserReq,
) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
