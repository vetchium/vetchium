package employer

import "github.com/vetchium/vetchium/typespec/common"

type ListHubUserEducationRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
