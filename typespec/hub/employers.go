package hub

type FilterEmployersRequest struct {
	Prefix string `query:"prefix"`
}

type HubEmployer struct {
	Domain    string `json:"domain"`
	Name      string `json:"name"`
	AsciiName string `json:"ascii_name"`
}

type FilterEmployersResponse struct {
	Employers []HubEmployer `json:"employers"`
}

type GetEmployerDetailsRequest struct {
	Domain string `json:"domain" validate:"required,validate_domain"`
}

type HubEmployerDetails struct {
	Name                   string `json:"name"`
	VerifiedEmployeesCount uint32 `json:"verified_employees_count"`

	// Anything below this line should be ignored if IsOnboarded is false
	IsOnboarded         bool   `json:"is_onboarded"`
	ActiveOpeningsCount uint32 `json:"active_openings_count"`
	IsFollowing         bool   `json:"is_following"`
}
