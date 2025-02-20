package hub

import "github.com/psankar/vetchi/typespec/common"

type ConnectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type UnlinkColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type MyColleagueApprovalsRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type MyColleagueSeeksRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type ApproveColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type RejectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
