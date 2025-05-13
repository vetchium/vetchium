package db

import (
	"context"

	"github.com/vetchium/vetchium/typespec/employer"
	"github.com/vetchium/vetchium/typespec/hub"
)

type AddPostRequest struct {
	Context context.Context
	PostID  string `json:"post_id"`
	hub.AddPostRequest
}

type GetPostRequest struct {
	Context context.Context
	PostID  string `json:"post_id"`
}

type AddEmployerPostRequest struct {
	Context context.Context
	PostID  string `json:"post_id"`
	employer.AddEmployerPostRequest
}

type UpdateEmployerPostRequest struct {
	Context context.Context
	PostID  string `json:"post_id"`
	employer.UpdateEmployerPostRequest
}

type ListEmployerPostsRequest struct {
	Context context.Context
	employer.ListEmployerPostsRequest
}
