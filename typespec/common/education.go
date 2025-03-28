package common

type Institute struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type Education struct {
	ID              string  `json:"id"`
	InstituteDomain string  `json:"institute_domain"`
	Degree          string  `json:"degree"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	Description     *string `json:"description"`
}
