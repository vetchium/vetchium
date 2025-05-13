package common

import "time"

type EmployerPost struct {
	ID            string    `json:"id"`
	Content       string    `json:"content"`
	Tags          []string  `json:"tags"`
	CompanyDomain string    `json:"company_domain"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
