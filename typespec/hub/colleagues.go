package hub

import "github.com/psankar/vetchi/typespec/common"

type HubUserShort struct {
	Handle   common.Handle `json:"handle"`
	Name     string        `json:"name"`
	ShortBio string        `json:"short_bio"`
}

type ConnectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type UnlinkColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type MyColleagueApprovalsRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int     `json:"limit"          validate:"min=0,max=100"`
}

type MyColleagueApprovals struct {
	Approvals     []HubUserShort `json:"approvals"`
	PaginationKey string         `json:"pagination_key,omitempty"`
}

type MyColleagueSeeksRequest struct {
	PaginationKey *string `json:"pagination_key"`
	Limit         int     `json:"limit"          validate:"min=1,max=100"`
}

type MyColleagueSeeks struct {
	Seeks         []HubUserShort `json:"seeks"`
	PaginationKey string         `json:"pagination_key,omitempty"`
}

type ApproveColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type RejectColleagueRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
