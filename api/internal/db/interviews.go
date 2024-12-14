package db

type AddInterviewersRequest struct {
	InterviewID string
	OrgUserIDs  []string
	Email       Email
}
