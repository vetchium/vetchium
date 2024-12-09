package vetchi

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
	Resume                 string `json:"resume"                    validate:"required"`
	Filename               string `json:"filename"                  validate:"required,max=256"`
	CoverLetter            string `json:"cover_letter"              validate:"omitempty,max=4096"`
}

type UpdateApplicationStateRequest struct {
	ID        string           `json:"id"         validate:"required"`
	FromState ApplicationState `json:"from_state" validate:"required"`
	ToState   ApplicationState `json:"to_state"   validate:"required"`
}
