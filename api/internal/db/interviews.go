package db

type AddInterviewerRequest struct {
	InterviewID string

	InterviewerEmailAddr string
	CandidacyComment     string

	InterviewerNotificationEmail Email
	WatcherNotificationEmail     Email
}

type RemoveInterviewerRequest struct {
	InterviewID      string
	CandidacyComment string

	RemovedInterviewerEmailAddr         string
	RemovedInterviewerEmailNotification Email

	// TODO: Add EmailNotification for watchers
}
