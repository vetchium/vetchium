package employer

type AddInterviewersRequest struct {
	InterviewID   string   `json:"interview_id"    validate:"required"`
	OrgUserEmails []string `json:"org_user_emails" validate:"required,dive,email"`
}

type RemoveInterviewerRequest struct {
	InterviewID  string `json:"interview_id"   validate:"required"`
	OrgUserEmail string `json:"org_user_email" validate:"required,email"`
}
