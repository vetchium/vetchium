package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetCandidaciesInfoRequest struct {
	RecruiterID   *string                `json:"recruiter_id"   validate:"omitempty"`
	State         *common.CandidacyState `json:"state"          validate:"omitempty"`
	PaginationKey *string                `json:"pagination_key" validate:"omitempty"`
	Limit         *int                   `json:"limit"          validate:"omitempty,max=40"`
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
	InPersonInterview    InterviewType = "IN_PERSON"
	VideoCallInterview   InterviewType = "VIDEO_CALL"
	TakeHomeInterview    InterviewType = "TAKE_HOME"
	UnspecifiedInterview InterviewType = "UNSPECIFIED"
)

type AddInterviewRequest struct {
	CandidacyID   string                `json:"candidacy_id"   validate:"required"`
	StartTime     time.Time             `json:"start_time"     validate:"required"`
	EndTime       time.Time             `json:"end_time"       validate:"required"`
	InterviewType InterviewType         `json:"interview_type" validate:"required"`
	Description   string                `json:"description"    validate:"required,max=2048"`
	Interviewers  []common.EmailAddress `json:"interviewers"   validate:"omitempty"`
}

type InterviewState string

const (
	ScheduledInterview         InterviewState = "SCHEDULED"
	CompletedInterview         InterviewState = "COMPLETED"
	CandidateWithdrewInterview InterviewState = "CANDIDATE_WITHDREW"
	EmployerWithdrewInterview  InterviewState = "EMPLOYER_WITHDREW"
)

type InterviewersDecision string

const (
	StrongYesDecision InterviewersDecision = "STRONG_YES"
	YesDecision       InterviewersDecision = "YES"
	NoDecision        InterviewersDecision = "NO"
	StrongNoDecision  InterviewersDecision = "STRONG_NO"
)

type Interview struct {
	InterviewID            string                `json:"interview_id"`
	InterviewState         InterviewState        `json:"interview_state"`
	StartTime              time.Time             `json:"start_time"`
	EndTime                time.Time             `json:"end_time"`
	InterviewType          InterviewType         `json:"interview_type"`
	Description            *string               `json:"description,omitempty"`
	Interviewers           []common.EmailAddress `json:"interviewers,omitempty"`
	InterviewersDecision   *InterviewersDecision `json:"interviewers_decision,omitempty"`
	InterviewersAssessment *string               `json:"interviewers_assessment,omitempty"`
	FeedbackToCandidate    *string               `json:"feedback_to_candidate,omitempty"`
	CreatedAt              time.Time             `json:"created_at"`
}

type GetInterviewsRequest struct {
	CandidacyID   *string `json:"candidacy_id"   validate:"omitempty"`
	OpeningID     *string `json:"opening_id"     validate:"omitempty"`
	PaginationKey *string `json:"pagination_key" validate:"omitempty"`
	Limit         *int    `json:"limit"          validate:"omitempty,max=100"`
}
