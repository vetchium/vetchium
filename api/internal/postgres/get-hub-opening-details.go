package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetHubOpeningDetails(
	ctx context.Context,
	req hub.GetHubOpeningDetailsRequest,
) (hub.HubOpeningDetails, error) {
	return hub.HubOpeningDetails{}, nil
}
