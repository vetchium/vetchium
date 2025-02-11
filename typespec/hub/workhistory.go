package hub

type AddWorkHistoryRequest struct {
	EmployerDomain string  `json:"employer_domain" validate:"required,validate_domain"`
	Title          string  `json:"title"           validate:"required,min=3,max=64"`
	StartDate      string  `json:"start_date"      validate:"required,datetime=2006-01-02"`
	EndDate        *string `json:"end_date"        validate:"omitempty,datetime=2006-01-02,date_after=StartDate"`
	Description    *string `json:"description"     validate:"omitempty,max=1024"`
}

type AddWorkHistoryResponse struct {
	WorkHistoryID string `json:"work_history_id"`
}

type WorkHistory struct {
	ID             string `json:"id"`
	EmployerDomain string `json:"employer_domain" validate:"required,validate_domain"`

	EmployerName *string `json:"employer_name"`
	EmployerLogo *string `json:"employer_logo"`

	Title       string  `json:"title"       validate:"required,min=3,max=64"`
	StartDate   string  `json:"start_date"  validate:"required,datetime=2006-01-02"`
	EndDate     *string `json:"end_date"    validate:"omitempty,datetime=2006-01-02,date_after=StartDate"`
	Description *string `json:"description" validate:"omitempty,max=1024"`
}

type UpdateWorkHistoryRequest struct {
	ID          string  `json:"id"          validate:"required,uuid"`
	Title       string  `json:"title"       validate:"required,min=3,max=64"`
	StartDate   string  `json:"start_date"  validate:"required,datetime=2006-01-02"`
	EndDate     *string `json:"end_date"    validate:"omitempty,datetime=2006-01-02,date_after=StartDate"`
	Description *string `json:"description" validate:"omitempty,max=1024"`
}

type ListWorkHistoryRequest struct {
	UserHandle *string `json:"user_handle" validate:"omitempty,min=3,max=64"`
}

type DeleteWorkHistoryRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}
