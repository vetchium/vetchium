package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetCandidaciesInfoRequest struct {
	RecruiterEmail *string `json:"recruiter_email" validate:"omitempty"`
	State          *string `json:"state"           validate:"omitempty"`
	PaginationKey  *string `json:"pagination_key"  validate:"omitempty"`
	Limit          int64   `json:"limit"           validate:"required,min=0,max=40"`
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

func (i InterviewType) IsValid() bool {
	switch i {
	case InPersonInterviewType,
		VideoCallInterviewType,
		TakeHomeInterviewType,
		UnspecifiedInterviewType:
		return true
	}
	return false
}

type AddInterviewRequest struct {
	CandidacyID   string                `json:"candidacy_id"   validate:"required"`
	StartTime     time.Time             `json:"start_time"     validate:"required"`
	EndTime       time.Time             `json:"end_time"       validate:"required"`
	InterviewType InterviewType         `json:"interview_type" validate:"required,validate_interview_type"`
	Description   string                `json:"description"    validate:"omitempty,max=2048"`
	Interviewers  []common.EmailAddress `json:"interviewers"   validate:"omitempty"`
}

type Interview struct {
	InterviewID          string                       `json:"interview_id"`
	InterviewState       common.InterviewState        `json:"interview_state"`
	StartTime            time.Time                    `json:"start_time"`
	EndTime              time.Time                    `json:"end_time"`
	InterviewType        InterviewType                `json:"interview_type"`
	Description          *string                      `json:"description"`
	Interviewers         []OrgUserTiny                `json:"interviewers"`
	InterviewersDecision *common.InterviewersDecision `json:"interviewers_decision"`
	Positives            *string                      `json:"positives"             validate:"omitempty,max=4096"`
	Negatives            *string                      `json:"negatives"             validate:"omitempty,max=4096"`
	OverallAssessment    *string                      `json:"overall_assessment"    validate:"omitempty,max=4096"`
	FeedbackToCandidate  string                       `json:"feedback_to_candidate" validate:"omitempty,max=4096"`
	FeedbackSubmittedBy  *OrgUserTiny                 `json:"feedback_submitted_by"`
	FeedbackSubmittedAt  *time.Time                   `json:"feedback_submitted_at"`
	CreatedAt            time.Time                    `json:"created_at"`
}
