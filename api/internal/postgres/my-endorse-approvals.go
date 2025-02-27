package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetMyEndorsementApprovals(
	ctx context.Context,
	req hub.MyEndorseApprovalsRequest,
) ([]hub.MyEndorseApproval, error) {
	return nil, nil
}
