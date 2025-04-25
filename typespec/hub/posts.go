package hub

import (
	"time"

	"github.com/vetchium/vetchium/typespec/common"
)

type AddPostRequest struct {
	Content string            `json:"content"  validate:"required,min=1,max=4096"`
	TagIDs  []common.VTagID   `json:"tag_ids"  validate:"max=3,dive,uuid"`
	NewTags []common.VTagName `json:"new_tags" validate:"max=3"`
}

type AddPostResponse struct {
	PostID string `json:"post_id"`
}

type Post struct {
	ID             string        `json:"id"`
	Content        string        `json:"content"`
	Tags           []string      `json:"tags"`
	AuthorName     string        `json:"author_name"`
	AuthorHandle   common.Handle `json:"author_handle"`
	CreatedAt      time.Time     `json:"created_at"`
	UpvotesCount   int32         `json:"upvotes_count"`
	DownvotesCount int32         `json:"downvotes_count"`
	Score          int32         `json:"score"`
	MeUpvoted      bool          `json:"me_upvoted"`
	MeDownvoted    bool          `json:"me_downvoted"`
	CanUpvote      bool          `json:"can_upvote"`
	CanDownvote    bool          `json:"can_downvote"`
	AmIAuthor      bool          `json:"am_i_author"`
}

type GetUserPostsRequest struct {
	Handle        *common.Handle `json:"handle"         validate:"omitempty,validate_handle"`
	PaginationKey *string        `json:"pagination_key" validate:"omitempty"`
	Limit         int            `json:"limit"          validate:"min=0,max=40"`
}

type GetUserPostsResponse struct {
	Posts         []Post `json:"posts"`
	PaginationKey string `json:"pagination_key"`
}

type FollowUserRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type UnfollowUserRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type GetFollowStatusRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type FollowStatus struct {
	IsFollowing bool `json:"is_following"`
	IsBlocked   bool `json:"is_blocked"`
	CanFollow   bool `json:"can_follow"`
}

type GetMyHomeTimelineRequest struct {
	PaginationKey *string `json:"pagination_key" validate:"omitempty"`
	Limit         int     `json:"limit"          validate:"min=0,max=40"`
}

type MyHomeTimeline struct {
	Posts         []Post `json:"posts"`
	PaginationKey string `json:"pagination_key"`
}

type GetPostDetailsRequest struct {
	PostID string `json:"post_id" validate:"required"`
}

type UpvoteUserPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}

type DownvoteUserPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}

type UnvoteUserPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}
