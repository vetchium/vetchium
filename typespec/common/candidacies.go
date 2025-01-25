package common

import (
	"time"
)

type GetCandidacyInfoRequest struct {
	CandidacyID string `json:"candidacy_id"`
}

type GetCandidacyCommentsRequest struct {
	CandidacyID string `json:"candidacy_id"`
}

type CommenterType string

const (
	CommenterTypeOrgUser CommenterType = "ORG_USER"
	CommenterTypeHubUser CommenterType = "HUB_USER"
)

type CandidacyComment struct {
	CommentID     string        `json:"comment_id"`
	CommenterName string        `json:"commenter_name"`
	CommenterType CommenterType `json:"commenter_type"`
	Content       string        `json:"content"`
	CreatedAt     time.Time     `json:"created_at"`
}
