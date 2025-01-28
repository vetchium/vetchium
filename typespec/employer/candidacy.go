package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type FilterCandidacyInfosRequest struct {
	OpeningID      *string `json:"opening_id"      validate:"omitempty"`
	RecruiterEmail *string `json:"recruiter_email" validate:"omitempty"`
	State          *string `json:"state"           validate:"omitempty"`

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

type AddInterviewRequest struct {
	CandidacyID       string               `json:"candidacy_id"       validate:"required"`
	StartTime         time.Time            `json:"start_time"         validate:"required"`
	EndTime           time.Time            `json:"end_time"           validate:"required"`
	InterviewType     common.InterviewType `json:"interview_type"     validate:"required,validate_interview_type"`
	Description       string               `json:"description"        validate:"omitempty,max=2048"`
	InterviewerEmails []string             `json:"interviewer_emails" validate:"omitempty,dive,email"`
}

type AddInterviewResponse struct {
	InterviewID string `json:"interview_id"`
}

type Interviewer struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	RSVPStatus common.RSVPStatus `json:"rsvp_status"`
}

type EmployerInterview struct {
	InterviewID          string                       `json:"interview_id"`
	InterviewState       common.InterviewState        `json:"interview_state"`
	StartTime            time.Time                    `json:"start_time"`
	EndTime              time.Time                    `json:"end_time"`
	InterviewType        common.InterviewType         `json:"interview_type"`
	Description          *string                      `json:"description"`
	CandidateName        string                       `json:"candidate_name"`
	CandidateHandle      string                       `json:"candidate_handle"`
	CandidateRSVPStatus  common.RSVPStatus            `json:"candidate_rsvp_status"`
	Interviewers         []Interviewer                `json:"interviewers"`
	InterviewersDecision *common.InterviewersDecision `json:"interviewers_decision"`
	Positives            *string                      `json:"positives"             validate:"omitempty,max=4096"`
	Negatives            *string                      `json:"negatives"             validate:"omitempty,max=4096"`
	OverallAssessment    *string                      `json:"overall_assessment"    validate:"omitempty,max=4096"`
	FeedbackToCandidate  *string                      `json:"feedback_to_candidate" validate:"omitempty,max=4096"`
	FeedbackSubmittedBy  *OrgUserTiny                 `json:"feedback_submitted_by"`
	FeedbackSubmittedAt  *time.Time                   `json:"feedback_submitted_at"`
	CreatedAt            time.Time                    `json:"created_at"`
}

type OfferToCandidateRequest struct {
	CandidacyID   string `json:"candidacy_id"   validate:"required"`
	OfferDocument string `json:"offer_document" validate:"omitempty"`
}
