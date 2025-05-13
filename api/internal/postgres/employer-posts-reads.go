package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func (p *PG) GetEmployerPost(
	ctx context.Context,
	postID string,
) (common.EmployerPost, error) {
	return common.EmployerPost{}, nil
}

func (p *PG) ListEmployerPosts(
	req db.ListEmployerPostsRequest,
) (employer.ListEmployerPostsResponse, error) {
	return employer.ListEmployerPostsResponse{}, nil
}
