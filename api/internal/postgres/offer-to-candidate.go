package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) OfferToCandidate(
	ctx context.Context,
	request employer.OfferToCandidateRequest,
) error {
	return nil
}
