package db

import (
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

const (
	// Corresponds to the comment_author_types enum in the database
	HubUserAuthorType = "HUB_USER"
	OrgUserAuthorType = "ORG_USER"
)

type Stakeholders struct {
	OpeningID    string
	OpeningTitle string
	OpeningState common.OpeningState
	OpeningType  common.OpeningType

	HiringManager employer.OrgUserShort
	Recruiter     employer.OrgUserShort

	Watchers []employer.OrgUserShort
}

type AddInterviewRequest struct {
	employer.AddInterviewRequest
	InterviewID string

	InterviewerNotificationEmail Email
	WatcherNotificationEmail     Email
	ApplicantNotificationEmail   Email
	CandidacyComment             string
}
