package postgres

import (
	"context"

	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) UpvoteUserPost(
	ctx context.Context,
	req hub.UpvoteUserPostRequest,
) error {
	return nil
}

func (pg *PG) DownvoteUserPost(
	ctx context.Context,
	req hub.DownvoteUserPostRequest,
) error {
	return nil
}

func (pg *PG) UnvoteUserPost(
	ctx context.Context,
	req hub.UnvoteUserPostRequest,
) error {
	return nil
}
