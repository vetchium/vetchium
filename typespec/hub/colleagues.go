package hub

import (
	"fmt"
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

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

type FilterColleaguesRequest struct {
	Prefix string `json:"prefix" validate:"required,min=1,max=12"`
	Limit  int    `json:"limit"  validate:"required,min=1,max=6"`
}

type EndorsementState string

const (
	SoughtEndorsement   EndorsementState = "SOUGHT_ENDORSEMENT"
	Endorsed            EndorsementState = "ENDORSED"
	DeclinedEndorsement EndorsementState = "DECLINED_ENDORSEMENT"
)

// Scan implements the Scanner interface for EndorsementState
func (e *EndorsementState) Scan(src any) error {
	if src == nil {
		*e = ""
		return nil
	}

	switch v := src.(type) {
	case []byte:
		*e = EndorsementState(v)
		return nil
	case string:
		*e = EndorsementState(v)
		return nil
	default:
		return fmt.Errorf("cannot convert %T to EndorsementState", src)
	}
}

// Value implements the driver Valuer interface for EndorsementState
func (e EndorsementState) Value() (interface{}, error) {
	return string(e), nil
}

type MyEndorseApprovalsRequest struct {
	PaginationKey *string            `json:"pagination_key"`
	Limit         int                `json:"limit"          validate:"min=0,max=100"`
	State         []EndorsementState `json:"state"          validate:"min=0,max=3"`
}

type MyEndorseApproval struct {
	ApplicationID        string           `json:"application_id"`
	ApplicantHandle      common.Handle    `json:"applicant_handle"`
	ApplicantName        string           `json:"applicant_name"`
	ApplicantShortBio    string           `json:"applicant_short_bio"`
	EmployerName         string           `json:"employer_name"`
	EmployerDomain       string           `json:"employer_domain"`
	OpeningTitle         string           `json:"opening_title"`
	OpeningURL           string           `json:"opening_url"`
	ApplicationStatus    string           `json:"application_status"`
	ApplicationCreatedAt time.Time        `json:"application_created_at"`
	EndorsementStatus    EndorsementState `json:"endorsement_status"`
}

type MyEndorseApprovalsResponse struct {
	Endorsements  []MyEndorseApproval `json:"endorsements"`
	PaginationKey string              `json:"pagination_key,omitempty"`
}

type EndorseApplicationRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}

type RejectEndorsementRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}
