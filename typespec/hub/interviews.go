package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetHubInterviewsByCandidacyRequest struct {
	CandidacyID string                  `json:"candidacy_id"`
	States      []common.InterviewState `json:"states"`
}

type HubInterview struct {
	InterviewID    string                `json:"interview_id"`
	InterviewState common.InterviewState `json:"interview_state"`
	StartTime      time.Time             `json:"start_time"`
	EndTime        time.Time             `json:"end_time"`
	InterviewType  common.InterviewType  `json:"interview_type"`
	Description    string                `json:"description"`
	Interviewers   []string              `json:"interviewers"`
}

type HubRSVPInterviewRequest struct {
	InterviewID string            `json:"interview_id"`
	RSVP        common.RSVPStatus `json:"rsvp"`
}
