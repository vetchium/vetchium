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

type PutAssesmentRequest struct {
	InterviewID         string                      `json:"interview_id"          validate:"required"`
	Decision            common.InterviewersDecision `json:"decision"              validate:"required"`
	Positives           string                      `json:"positives"             validate:"omitempty,max=4096"`
	Negatives           string                      `json:"negatives"             validate:"omitempty,max=4096"`
	OverallAssessment   string                      `json:"overall_assessment"    validate:"omitempty,max=4096"`
	FeedbackToCandidate string                      `json:"feedback_to_candidate" validate:"omitempty,max=4096"`
}
