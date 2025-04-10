package hub

import "github.com/vetchium/vetchium/typespec/common"

type AddPostRequest struct {
	Content string        `json:"content"  validate:"required,min=1,max=4096"`
	TagIDs  []string      `json:"tag_ids"  validate:"max=3,dive,uuid"`
	NewTags []common.VTag `json:"new_tags" validate:"max=3"`
}

type AddPostResponse struct {
	PostID string `json:"post_id"`
}

type GetTimelineRequest struct {
	TimelineID    *string `json:"timeline_id"`
	PaginationKey *string `json:"pagination_key"`
	Limit         int     `json:"limit"`
}

type GetTimelineResponse struct {
	Posts         []common.Post `json:"posts"`
	PaginationKey string        `json:"pagination_key"`
}
