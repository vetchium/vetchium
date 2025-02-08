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
