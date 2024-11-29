package vetchi

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
	Resume                 string `json:"resume"                    validate:"required"`
	Filename               string `json:"filename"                  validate:"required,max=256"`
	CoverLetter            string `json:"cover_letter"              validate:"omitempty,max=4096"`
}

type ApplicationState string

const (
	AppliedAppState     ApplicationState = "APPLIED"
	RejectedAppState    ApplicationState = "REJECTED"
	ShortlistedAppState ApplicationState = "SHORTLISTED"
	WithdrawnAppState   ApplicationState = "WITHDRAWN"
	ExpiredAppState     ApplicationState = "EXPIRED"
)
