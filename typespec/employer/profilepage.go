package employer

type GetHubUserBioRequest struct {
	Handle string `json:"handle" validate:"required"`
}

type EmployerViewBio struct {
	Handle              string   `json:"handle"`
	FullName            string   `json:"full_name"`
	ShortBio            string   `json:"short_bio"`
	LongBio             string   `json:"long_bio"`
	VerifiedMailDomains []string `json:"verified_mail_domains,omitempty"`
}
