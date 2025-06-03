package hub

import (
	"time"

	"github.com/vetchium/vetchium/typespec/common"
)

type IncognitoPost struct {
	IncognitoPostID string        `json:"incognito_post_id"`
	Content         string        `json:"content"`
	Tags            []common.VTag `json:"tags"`
	CreatedAt       time.Time     `json:"created_at"`
	Upvotes         int           `json:"upvotes"`
	Downvotes       int           `json:"downvotes"`
	IsCreatedByMe   bool          `json:"is_created_by_me"`
	IsDeleted       bool          `json:"is_deleted"`
}

type AddIncognitoPostRequest struct {
	Content string          `json:"content" validate:"required,min=1,max=1024"`
	TagIDs  []common.VTagID `json:"tag_ids" validate:"max=3"`
}

type AddIncognitoPostResponse struct {
	IncognitoPostID string `json:"incognito_post_id"`
}

type IncognitoPostComment struct {
	CommentID     string    `json:"comment_id"`
	Content       string    `json:"content"`
	InReplyTo     *string   `json:"in_reply_to"`
	CreatedAt     time.Time `json:"created_at"`
	Upvotes       int       `json:"upvotes"`
	Downvotes     int       `json:"downvotes"`
	IsCreatedByMe bool      `json:"is_created_by_me"`
	IsDeleted     bool      `json:"is_deleted"`
	MyVote        string    `json:"my_vote"`
	Depth         int       `json:"depth"`
}

type AddIncognitoPostCommentRequest struct {
	IncognitoPostID string  `json:"incognito_post_id" validate:"required"`
	Content         string  `json:"content"           validate:"required,min=1,max=512"`
	InReplyTo       *string `json:"in_reply_to"`
}

type GetIncognitoPostCommentsRequest struct {
	IncognitoPostID    string  `json:"incognito_post_id"    validate:"required"`
	PaginationKey      *string `json:"pagination_key"`
	Limit              int     `json:"limit"                validate:"required,min=1,max=100"`
	InReplyTo          *string `json:"in_reply_to"`
	IncludeNestedDepth int     `json:"include_nested_depth" validate:"required,min=0,max=5"`
}

type GetIncognitoPostCommentsResponse struct {
	Comments      []IncognitoPostComment `json:"comments"`
	PaginationKey *string                `json:"pagination_key"`
	HasMore       bool                   `json:"has_more"`
	TotalCount    int                    `json:"total_count"`
}

type DeleteIncognitoPostCommentRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
	CommentID       string `json:"comment_id"        validate:"required"`
}

type UpvoteIncognitoPostCommentRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
	CommentID       string `json:"comment_id"        validate:"required"`
}

type DownvoteIncognitoPostCommentRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
	CommentID       string `json:"comment_id"        validate:"required"`
}

type UnvoteIncognitoPostCommentRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
	CommentID       string `json:"comment_id"        validate:"required"`
}

type GetIncognitoPostRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
}

type DeleteIncognitoPostRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
}
