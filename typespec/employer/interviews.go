package employer

import "github.com/psankar/vetchi/typespec/common"

type AddInterviewerRequest struct {
	InterviewID  string `json:"interview_id"   validate:"required"`
	OrgUserEmail string `json:"org_user_email" validate:"required,email"`
}

type RemoveInterviewerRequest struct {
	InterviewID  string `json:"interview_id"   validate:"required"`
	OrgUserEmail string `json:"org_user_email" validate:"required,email"`
}

type GetInterviewsByOpeningRequest struct {
	OpeningID     string                  `json:"opening_id"     validate:"required"`
	States        []common.InterviewState `json:"states"         validate:"omitempty"`
	PaginationKey string                  `json:"pagination_key" validate:"omitempty"`
	Limit         int64                   `json:"limit"          validate:"required,min=0,max=100"`
}

type GetInterviewsByCandidacyRequest struct {
	CandidacyID string                  `json:"candidacy_id" validate:"required"`
	States      []common.InterviewState `json:"states"       validate:"omitempty"`
}
