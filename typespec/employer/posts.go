package employer

import "github.com/vetchium/vetchium/typespec/common"

type AddEmployerPostRequest struct {
	Content string          `json:"content" validate:"required,min=1,max=4096"`
	TagIDs  []common.VTagID `json:"tag_ids" validate:"max=3"`
}

type AddEmployerPostResponse struct {
	PostID string `json:"post_id"`
}

type UpdateEmployerPostRequest struct {
	PostID  string          `json:"post_id"`
	Content string          `json:"content" validate:"required,min=1,max=4096"`
	TagIDs  []common.VTagID `json:"tags"    validate:"max=3"`
}

type DeleteEmployerPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}

type ListEmployerPostsRequest struct {
	PaginationKey string `json:"pagination_key"`
	Limit         int    `json:"limit"          validate:"min=0,max=40"`
}

type ListEmployerPostsResponse struct {
	Posts         []common.EmployerPost `json:"posts"`
	PaginationKey string                `json:"pagination_key"`
}

type GetEmployerPostRequest struct {
	PostID string `json:"post_id" validate:"required"`
}
