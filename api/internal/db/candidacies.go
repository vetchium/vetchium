package db

import (
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

const (
	// Corresponds to the comment_author_types enum in the database
	HubUserAuthorType = "HUB_USER"
	OrgUserAuthorType = "ORG_USER"
)

type WatchersInfo struct {
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

	// Names of interviewers, used for applicant notification
	InterviewerNames []string
}
