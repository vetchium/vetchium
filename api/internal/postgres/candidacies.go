package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetCandidaciesInfo(
	ctx context.Context,
	getCandidaciesInfoReq employer.GetCandidaciesInfoRequest,
) ([]employer.Candidacy, error) {
	return nil, nil
}
