package hub

import "github.com/psankar/vetchi/typespec/common"

type HubRSVPInterviewRequest struct {
	InterviewID string            `json:"interview_id"`
	RSVP        common.RSVPStatus `json:"rsvp"`
}
