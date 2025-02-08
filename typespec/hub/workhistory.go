package hub

type AddWorkHistoryRequest struct {
	EmployerDomain string  `json:"employer_domain"`
	Title          string  `json:"title"`
	StartDate      string  `json:"start_date"`
	EndDate        *string `json:"end_date"`
	Description    *string `json:"description"     validate:"max=1024"`
}

type AddWorkHistoryResponse struct {
	WorkHistoryID string `json:"work_history_id"`
}

type WorkHistory struct {
	ID             string `json:"id"`
	EmployerDomain string `json:"employer_domain"`

	EmployerName *string `json:"employer_name"`
	EmployerLogo *string `json:"employer_logo"`

	Title       string  `json:"title"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description *string `json:"description"`
}

type UpdateWorkHistoryRequest struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description *string `json:"description" validate:"max=1024"`
}

type ListWorkHistoryRequest struct {
	UserHandle *string `json:"user_handle"`
}

type DeleteWorkHistoryRequest struct {
	ID string `json:"id"`
}
