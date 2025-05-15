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
	VerifiedEmployeesCount int    `json:"verified_employees_count"`
	ActiveOpeningsCount    int    `json:"active_openings_count"`
}
