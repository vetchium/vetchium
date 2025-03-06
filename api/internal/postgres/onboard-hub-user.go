package postgres

import (
	"context"

	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) OnboardHubUser(
	ctx context.Context,
	onboardHubUserReq hub.OnboardHubUserRequest,
) (hub.OnboardHubUserResponse, error) {
	return hub.OnboardHubUserResponse{}, nil
}
