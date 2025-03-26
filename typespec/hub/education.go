package hub

import "github.com/psankar/vetchi/typespec/common"

type AddEducationRequest struct {
	InstituteDomain common.Domain `json:"institute_domain" validate:"required,validate_domain"`
	Degree          string        `json:"degree"           validate:"required,min=3,max=64"`
	StartDate       *string       `json:"start_date"       validate:"omitempty,validate_date,no_future_date,required_with=EndDate"`
	EndDate         *string       `json:"end_date"         validate:"omitempty,validate_date,date_after=StartDate"`
	Description     *string       `json:"description"      validate:"omitempty,max=1024"`
}

type AddEducationResponse struct {
	EducationID string `json:"education_id"`
}

type Institute struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

type FilterInstitutesRequest struct {
	Prefix string `json:"prefix"`
}

type DeleteEducationRequest struct {
	EducationID string `json:"education_id"`
}

type Education struct {
	ID              string  `json:"id"`
	InstituteDomain string  `json:"institute_domain"`
	Degree          string  `json:"degree"`
	StartDate       *string `json:"start_date"`
	EndDate         *string `json:"end_date"`
	Description     *string `json:"description"`
}

type ListEducationRequest struct {
	UserHandle *string `json:"user_handle"`
}
