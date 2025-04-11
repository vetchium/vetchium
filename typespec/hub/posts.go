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
	Limit         int            `json:"limit"          validate:"min=1,max=40"`
}

type GetUserPostsResponse struct {
	Posts         []common.Post `json:"posts"`
	PaginationKey string        `json:"pagination_key"`
}
