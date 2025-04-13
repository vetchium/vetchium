package postgres

import (
	"context"

	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) FollowUser(ctx context.Context, handle string) error {
	pg.log.Inf("Entered PG FollowUser", "handle", handle)
	return nil
}

func (pg *PG) UnfollowUser(ctx context.Context, handle string) error {
	pg.log.Inf("Entered PG UnfollowUser", "handle", handle)
	return nil
}

func (pg *PG) GetFollowStatus(
	ctx context.Context,
	handle string,
) (hub.FollowStatus, error) {
	pg.log.Inf("Entered PG GetFollowStatus", "handle", handle)
	return hub.FollowStatus{}, nil
}
