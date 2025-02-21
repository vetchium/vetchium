package hub

import "github.com/psankar/vetchi/typespec/common"

type ConnectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type UnlinkColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type MyColleagueApprovalsRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int     `json:"limit"          validate:"min=1,max=100"`
}

type MyColleagueSeeksRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int     `json:"limit"          validate:"min=1,max=100"`
}

type ApproveColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type RejectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
