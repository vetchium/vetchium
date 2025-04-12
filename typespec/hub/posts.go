package hub

import "github.com/vetchium/vetchium/typespec/common"

type AddPostRequest struct {
	Content string            `json:"content"  validate:"required,min=1,max=4096"`
	TagIDs  []common.VTagID   `json:"tag_ids"  validate:"max=3,dive,uuid"`
	NewTags []common.VTagName `json:"new_tags" validate:"max=3"`
}

type AddPostResponse struct {
	PostID string `json:"post_id"`
}

type GetUserPostsRequest struct {
	Handle        *common.Handle `json:"handle"         validate:"omitempty,validate_handle"`
	PaginationKey *string        `json:"pagination_key" validate:"omitempty"`
	Limit         int            `json:"limit"          validate:"min=0,max=40"`
}

type GetUserPostsResponse struct {
	Posts         []common.Post `json:"posts"`
	PaginationKey string        `json:"pagination_key"`
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
