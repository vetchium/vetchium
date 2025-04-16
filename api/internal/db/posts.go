package db

import (
	"context"

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
