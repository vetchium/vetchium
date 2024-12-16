package employer

type AddInterviewerRequest struct {
	InterviewID  string `json:"interview_id"   validate:"required"`
	OrgUserEmail string `json:"org_user_email" validate:"required,email"`
}

type RemoveInterviewerRequest struct {
	InterviewID  string `json:"interview_id"   validate:"required"`
	OrgUserEmail string `json:"org_user_email" validate:"required,email"`
}
