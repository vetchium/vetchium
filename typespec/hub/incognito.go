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
	UpvotesCount    int32         `json:"upvotes_count"`
	DownvotesCount  int32         `json:"downvotes_count"`
	Score           int32         `json:"score"`
	MeUpvoted       bool          `json:"me_upvoted"`
	MeDownvoted     bool          `json:"me_downvoted"`
	CanUpvote       bool          `json:"can_upvote"`
	CanDownvote     bool          `json:"can_downvote"`
	IsCreatedByMe   bool          `json:"is_created_by_me"`
	IsDeleted       bool          `json:"is_deleted"`
}

type AddIncognitoPostRequest struct {
	Content string          `json:"content" validate:"required,min=1,max=1024"`
	TagIDs  []common.VTagID `json:"tag_ids" validate:"min=1,max=3"`
}

type AddIncognitoPostResponse struct {
	IncognitoPostID string `json:"incognito_post_id"`
}

type IncognitoPostComment struct {
	CommentID      string    `json:"comment_id"`
	Content        string    `json:"content"`
	InReplyTo      *string   `json:"in_reply_to"`
	CreatedAt      time.Time `json:"created_at"`
	UpvotesCount   int32     `json:"upvotes_count"`
	DownvotesCount int32     `json:"downvotes_count"`
	Score          int32     `json:"score"`
	MeUpvoted      bool      `json:"me_upvoted"`
	MeDownvoted    bool      `json:"me_downvoted"`
	CanUpvote      bool      `json:"can_upvote"`
	CanDownvote    bool      `json:"can_downvote"`
	IsCreatedByMe  bool      `json:"is_created_by_me"`
	IsDeleted      bool      `json:"is_deleted"`
	Depth          int32     `json:"depth"`
	RepliesCount   int32     `json:"replies_count"`
}

type AddIncognitoPostCommentRequest struct {
	IncognitoPostID string  `json:"incognito_post_id" validate:"required"`
	Content         string  `json:"content"           validate:"required,min=1,max=512"`
	InReplyTo       *string `json:"in_reply_to"`
}

type AddIncognitoPostCommentResponse struct {
	IncognitoPostID string `json:"incognito_post_id"`
	CommentID       string `json:"comment_id"`
}

type IncognitoPostCommentSortBy string

const (
	IncognitoPostCommentSortByTop IncognitoPostCommentSortBy = "top"
	IncognitoPostCommentSortByNew IncognitoPostCommentSortBy = "new"
	IncognitoPostCommentSortByOld IncognitoPostCommentSortBy = "old"
)

type GetIncognitoPostCommentsRequest struct {
	IncognitoPostID         string                     `json:"incognito_post_id"          validate:"required"`
	SortBy                  IncognitoPostCommentSortBy `json:"sort_by"`
	PaginationKey           *string                    `json:"pagination_key"`
	Limit                   int32                      `json:"limit"                      validate:"min=1,max=50"`
	DirectRepliesPerComment int32                      `json:"direct_replies_per_comment" validate:"min=0,max=10"`
}

type GetIncognitoPostCommentsResponse struct {
	Comments           []IncognitoPostComment `json:"comments"`
	PaginationKey      string                 `json:"pagination_key"`
	TotalCommentsCount int32                  `json:"total_comments_count"`
}

type GetCommentRepliesRequest struct {
	IncognitoPostID string  `json:"incognito_post_id" validate:"required"`
	ParentCommentID string  `json:"parent_comment_id" validate:"required"`
	PaginationKey   *string `json:"pagination_key"`
	Limit           int32   `json:"limit"             validate:"min=1,max=100"`
	DirectOnly      bool    `json:"direct_only"`
	MaxDepth        int32   `json:"max_depth"         validate:"min=1,max=3"`
}

type GetCommentRepliesResponse struct {
	Replies           []IncognitoPostComment `json:"replies"`
	PaginationKey     string                 `json:"pagination_key"`
	TotalRepliesCount int32                  `json:"total_replies_count"`
	ParentCommentID   string                 `json:"parent_comment_id"`
}

type GetIncognitoPostCommentPermalinkRequest struct {
	IncognitoPostID      string `json:"incognito_post_id"      validate:"required"`
	CommentID            string `json:"comment_id"             validate:"required"`
	ContextSiblingsCount int32  `json:"context_siblings_count" validate:"min=0,max=10"`
	ContextRepliesCount  int32  `json:"context_replies_count"  validate:"min=0,max=20"`
}

type GetIncognitoPostCommentPermalinkResponse struct {
	Comments        []IncognitoPostComment `json:"comments"`
	TargetCommentID string                 `json:"target_comment_id"`
	BreadcrumbPath  []string               `json:"breadcrumb_path"`
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

type UpvoteIncognitoPostRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
}

type DownvoteIncognitoPostRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
}

type UnvoteIncognitoPostRequest struct {
	IncognitoPostID string `json:"incognito_post_id" validate:"required"`
}

type MyIncognitoPostComment struct {
	CommentID          string        `json:"comment_id"`
	Content            string        `json:"content"`
	InReplyTo          *string       `json:"in_reply_to"`
	CreatedAt          time.Time     `json:"created_at"`
	UpvotesCount       int32         `json:"upvotes_count"`
	DownvotesCount     int32         `json:"downvotes_count"`
	Score              int32         `json:"score"`
	MeUpvoted          bool          `json:"me_upvoted"`
	MeDownvoted        bool          `json:"me_downvoted"`
	IsDeleted          bool          `json:"is_deleted"`
	Depth              int32         `json:"depth"`
	IncognitoPostID    string        `json:"incognito_post_id"`
	PostContentPreview string        `json:"post_content_preview"`
	PostTags           []common.VTag `json:"post_tags"`
}

type GetMyIncognitoPostCommentsRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int32   `json:"limit"          validate:"min=1,max=40"`
}

type GetMyIncognitoPostCommentsResponse struct {
	Comments      []MyIncognitoPostComment `json:"comments"`
	PaginationKey string                   `json:"pagination_key"`
}

type IncognitoPostTimeFilter string

const (
	IncognitoPostTimeFilterPast24Hours IncognitoPostTimeFilter = "past_24_hours"
	IncognitoPostTimeFilterPastWeek    IncognitoPostTimeFilter = "past_week"
	IncognitoPostTimeFilterPastMonth   IncognitoPostTimeFilter = "past_month"
	IncognitoPostTimeFilterPastYear    IncognitoPostTimeFilter = "past_year"
)

type GetIncognitoPostsRequest struct {
	TagID         common.VTagID            `json:"tag_id"         validate:"required"`
	TimeFilter    *IncognitoPostTimeFilter `json:"time_filter"`
	Limit         int32                    `json:"limit"          validate:"min=1,max=100"`
	PaginationKey *string                  `json:"pagination_key"`
}

type IncognitoPostSummary struct {
	IncognitoPostID string        `json:"incognito_post_id"`
	Content         string        `json:"content"`
	Tags            []common.VTag `json:"tags"`
	CreatedAt       time.Time     `json:"created_at"`
	UpvotesCount    int32         `json:"upvotes_count"`
	DownvotesCount  int32         `json:"downvotes_count"`
	Score           int32         `json:"score"`
	MeUpvoted       bool          `json:"me_upvoted"`
	MeDownvoted     bool          `json:"me_downvoted"`
	CanUpvote       bool          `json:"can_upvote"`
	CanDownvote     bool          `json:"can_downvote"`
	CommentsCount   int32         `json:"comments_count"`
	IsCreatedByMe   bool          `json:"is_created_by_me"`
	IsDeleted       bool          `json:"is_deleted"`
}

type GetIncognitoPostsResponse struct {
	Posts         []IncognitoPostSummary `json:"posts"`
	PaginationKey string                 `json:"pagination_key"`
}

type GetMyIncognitoPostsRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int32   `json:"limit"          validate:"min=1,max=40"`
}

type GetMyIncognitoPostsResponse struct {
	Posts         []IncognitoPostSummary `json:"posts"`
	PaginationKey string                 `json:"pagination_key"`
}
