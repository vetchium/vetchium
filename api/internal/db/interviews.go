package db

type AddInterviewersRequest struct {
	InterviewID  string
	Interviewers []string
	Email        Email
}
