package vetchi

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company"`
	CompanyDomain          string `json:"company_domain"`
	Resume                 string `json:"resume"`
	Filename               string `json:"filename"                  validate:"max=256"`
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
