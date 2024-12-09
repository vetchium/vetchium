package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetCandidaciesInfoRequest struct {
	RecruiterID   *string `json:"recruiter_id"   validate:"omitempty"`
	State         *string `json:"state"          validate:"omitempty"`
	PaginationKey *string `json:"pagination_key" validate:"omitempty"`
	Limit         int64   `json:"limit"          validate:"required,min=0,max=40"`
}

type Candidacy struct {
	CandidacyID        string                `json:"candidacy_id"`
	OpeningID          string                `json:"opening_id"`
	OpeningTitle       string                `json:"opening_title"`
	OpeningDescription string                `json:"opening_description"`
	CandidacyState     common.CandidacyState `json:"candidacy_state"`
	ApplicantName      string                `json:"applicant_name"`
	ApplicantHandle    string                `json:"applicant_handle"`
}

type AddEmployerCandidacyCommentRequest struct {
	CandidacyID string `json:"candidacy_id" validate:"required"`
	Comment     string `json:"comment"      validate:"required,max=2048"`
}

type InterviewType string

const (
	InPersonInterviewType    InterviewType = "IN_PERSON"
	VideoCallInterviewType   InterviewType = "VIDEO_CALL"
	TakeHomeInterviewType    InterviewType = "TAKE_HOME"
	UnspecifiedInterviewType InterviewType = "UNSPECIFIED"
)

type AddInterviewRequest struct {
	CandidacyID   string                `json:"candidacy_id"   validate:"required"`
	StartTime     time.Time             `json:"start_time"     validate:"required"`
	EndTime       time.Time             `json:"end_time"       validate:"required"`
	InterviewType InterviewType         `json:"interview_type" validate:"required"`
	Description   string                `json:"description"    validate:"omitempty,max=2048"`
	Interviewers  []common.EmailAddress `json:"interviewers"   validate:"omitempty"`
}

type Interview struct {
	InterviewID            string                      `json:"interview_id"`
	InterviewState         common.InterviewState       `json:"interview_state"`
	StartTime              time.Time                   `json:"start_time"`
	EndTime                time.Time                   `json:"end_time"`
	InterviewType          InterviewType               `json:"interview_type"`
	Description            string                      `json:"description"`
	Interviewers           []common.EmailAddress       `json:"interviewers"`
	InterviewersDecision   common.InterviewersDecision `json:"interviewers_decision"`
	InterviewersAssessment string                      `json:"interviewers_assessment" validate:"omitempty,max=4096"`
	FeedbackToCandidate    string                      `json:"feedback_to_candidate"   validate:"omitempty,max=4096"`
	CreatedAt              time.Time                   `json:"created_at"`
}

type GetInterviewsRequest struct {
	CandidacyID   *string `json:"candidacy_id"   validate:"omitempty"`
	OpeningID     *string `json:"opening_id"     validate:"omitempty"`
	PaginationKey *string `json:"pagination_key" validate:"omitempty"`
	Limit         int64   `json:"limit"          validate:"required,min=0,max=100"`
}
