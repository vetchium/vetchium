package db

import "github.com/vetchium/vetchium/typespec/hub"

type AddPostCommentRequest struct {
	CommentID string
	hub.AddPostCommentRequest
}
