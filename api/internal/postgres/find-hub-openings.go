package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) FindHubOpenings(
	ctx context.Context,
	req *vetchi.FindHubOpeningsRequest,
) ([]vetchi.HubOpening, error) {
	return nil, nil
}
