package db

type AddInterviewersRequest struct {
	InterviewID  string
	Interviewers []string
	Email        Email

	// TODO: Add EmailNotification for watchers
	// TODO: Add CandidacyComment
}

type RemoveInterviewerRequest struct {
	InterviewID      string
	CandidacyComment string

	RemovedInterviewerEmailAddr         string
	RemovedInterviewerEmailNotification Email

	// TODO: Add EmailNotification for watchers
}
