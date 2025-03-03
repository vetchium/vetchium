package employer

type GetHubUserBioRequest struct {
	Handle string `json:"handle" validate:"required"`
}

type EmployerWorkHistory struct {
	ID             string  `json:"id"`
	EmployerDomain string  `json:"employer_domain"`
	EmployerName   *string `json:"employer_name"`
	Title          string  `json:"title"`
	StartDate      string  `json:"start_date"`
	EndDate        *string `json:"end_date"`
	Description    *string `json:"description"`
}

type EmployerViewBio struct {
	Handle              string                `json:"handle"`
	FullName            string                `json:"full_name"`
	ShortBio            string                `json:"short_bio"`
	LongBio             string                `json:"long_bio"`
	VerifiedMailDomains []string              `json:"verified_mail_domains,omitempty"`
	WorkHistory         []EmployerWorkHistory `json:"work_history"`
}
