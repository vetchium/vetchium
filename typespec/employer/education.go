package employer

import "github.com/psankar/vetchi/typespec/common"

type ListHubUserEducationRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
