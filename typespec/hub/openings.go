package hub

import "github.com/psankar/vetchi/typespec/common"

type ExperienceRange struct {
	YoeMin int `json:"yoe_min" validate:"min=0,max=99"`
	YoeMax int `json:"yoe_max" validate:"min=1,max=100"`
}

type SalaryRange struct {
	Currency common.Currency `json:"currency" validate:"validate_currency"`
	Min      float64         `json:"min"      validate:"min=0"`
	Max      float64         `json:"max"      validate:"min=1"`
}

type LocationFilter struct {
	CountryCode common.CountryCode `json:"country_code" validate:"validate_country_code"`
	City        string             `json:"city"         validate:"min=3,max=32"`
}

// Helper function to convert string to CountryCode pointer
func CountryCodePtr(c string) *common.CountryCode {
	cc := common.CountryCode(c)
	return &cc
}

type FindHubOpeningsRequest struct {
	CountryCode common.CountryCode `json:"country_code" validate:"required,validate_country_code"`
	Cities      []string           `json:"cities"       validate:"dive,omitempty"`

	OpeningTypes    []common.OpeningType `json:"opening_types"    validate:"dive,omitempty,validate_opening_type"`
	CompanyDomains  []string             `json:"company_domains"  validate:"dive,omitempty,validate_domain"`
	ExperienceRange *ExperienceRange     `json:"experience_range" validate:"omitempty"`
	SalaryRange     *SalaryRange         `json:"salary_range"     validate:"omitempty"`

	MinEducationLevel  *common.EducationLevel `json:"min_education_level"  validate:"omitempty,validate_education_level"`
	RemoteTimezones    []common.TimeZone      `json:"remote_timezones"     validate:"dive,omitempty,validate_timezone"`
	RemoteCountryCodes []common.CountryCode   `json:"remote_country_codes" validate:"dive,omitempty,validate_country_code"`
	PaginationKey      int64                  `json:"pagination_key"`
	Limit              int64                  `json:"limit"                validate:"min=1,max=100"`
}

type HubOpening struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company"`
	CompanyDomain          string `json:"company_domain"`
	CompanyName            string `json:"company_name"`
	JobTitle               string `json:"job_title"`
	JD                     string `json:"jd"`
	PaginationKey          int64  `json:"pagination_key"`
}

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
	Resume                 string `json:"resume"                    validate:"required"`
	Filename               string `json:"filename"                  validate:"required,max=256"`
	CoverLetter            string `json:"cover_letter"              validate:"omitempty,max=4096"`
}
