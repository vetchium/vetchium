package hub

import (
	"time"

	"github.com/vetchium/vetchium/typespec/common"
)

type AddPostCommentRequest struct {
	PostID  string `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=4096"`
}

type AddPostCommentResponse struct {
	PostID    string `json:"post_id"`
	CommentID string `json:"comment_id"`
}

type GetPostCommentsRequest struct {
	PostID        string `json:"post_id"        validate:"required"`
	PaginationKey string `json:"pagination_key"`
	Limit         int    `json:"limit"          validate:"min=0,max=40"`
}

type PostComment struct {
	ID           string        `json:"id"`
	Content      string        `json:"content"`
	AuthorName   string        `json:"author_name"`
	AuthorHandle common.Handle `json:"author_handle"`
	CreatedAt    time.Time     `json:"created_at"`
}

type DisablePostCommentsRequest struct {
	PostID                 string `json:"post_id"                  validate:"required"`
	DeleteExistingComments bool   `json:"delete_existing_comments"`
}

type EnablePostCommentsRequest struct {
	PostID string `json:"post_id" validate:"required"`
}

type DeletePostCommentRequest struct {
	PostID    string `json:"post_id"    validate:"required"`
	CommentID string `json:"comment_id" validate:"required"`
}

type DeleteMyCommentRequest struct {
	PostID    string `json:"post_id"    validate:"required"`
	CommentID string `json:"comment_id" validate:"required"`
}
