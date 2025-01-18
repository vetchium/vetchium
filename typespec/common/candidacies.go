package common

import (
	"time"
)

type GetCandidacyInfoRequest struct {
	CandidacyID string `json:"candidacy_id"`
}

type Candidacy struct {
	CandidacyID        string         `json:"candidacy_id"`
	OpeningID          string         `json:"opening_id"`
	OpeningTitle       string         `json:"opening_title"`
	OpeningDescription string         `json:"opening_description"`
	CandidacyState     CandidacyState `json:"candidacy_state"`
	ApplicantName      string         `json:"applicant_name"`
	ApplicantHandle    string         `json:"applicant_handle"`
}

type GetCandidacyCommentsRequest struct {
	CandidacyID string `json:"candidacy_id"`
}

type CandidacyComment struct {
	CommentID     string    `json:"comment_id"`
	CommenterName string    `json:"commenter_name"`
	CommenterType string    `json:"commenter_type"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}
