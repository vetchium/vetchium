package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetHubInterviewsByCandidacyRequest struct {
	CandidacyID string                  `json:"candidacy_id"`
	States      []common.InterviewState `json:"states"`
}

type HubInterviewer struct {
	Name       string            `json:"name"`
	RSVPStatus common.RSVPStatus `json:"rsvp_status"`
}

type HubInterview struct {
	InterviewID    string                `json:"interview_id"`
	InterviewState common.InterviewState `json:"interview_state"`
	StartTime      time.Time             `json:"start_time"`
	EndTime        time.Time             `json:"end_time"`
	InterviewType  common.InterviewType  `json:"interview_type"`
	Description    string                `json:"description"`
	CandidateRSVP  common.RSVPStatus     `json:"candidate_rsvp_status"`
	Interviewers   []HubInterviewer      `json:"interviewers"`
}

type HubRSVPInterviewRequest struct {
	InterviewID string            `json:"interview_id"`
	RSVP        common.RSVPStatus `json:"rsvp"`
}
