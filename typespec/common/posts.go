package common

import "time"

type EmployerPost struct {
	ID                 string    `json:"id"`
	Content            string    `json:"content"`
	Tags               []string  `json:"tags"`
	EmployerName       string    `json:"employer_name"`
	EmployerDomainName string    `json:"employer_domain_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
