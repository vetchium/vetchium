package common

type EmployerPost struct {
	ID            string   `json:"id"`
	Content       string   `json:"content"`
	Tags          []string `json:"tags"`
	CompanyDomain string   `json:"company_domain"`
	CreatedAt     string   `json:"created_at"`
}
