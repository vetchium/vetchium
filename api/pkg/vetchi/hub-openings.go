package vetchi

type ExperienceRange struct {
	YoeMin int `json:"yoe_min" validate:"min=0,max=99"`
	YoeMax int `json:"yoe_max" validate:"min=1,max=100"`
}

type SalaryRange struct {
	Currency Currency `json:"currency" validate:"validate_currency"`
	Min      float64  `json:"min"      validate:"min=0"`
	Max      float64  `json:"max"      validate:"min=1"`
}

type LocationFilter struct {
	CountryCode CountryCode `json:"country_code" validate:"validate_country_code"`
	City        string      `json:"city"         validate:"min=3,max=32"`
}

// Helper function to convert string to CountryCode pointer
func CountryCodePtr(c string) *CountryCode {
	cc := CountryCode(c)
	return &cc
}

type FindHubOpeningsRequest struct {
	CountryCode *CountryCode `json:"country_code" validate:"omitempty,validate_country_code"`
	Cities      []string     `json:"cities"       validate:"dive,omitempty,min=1,max=4"`

	OpeningTypes    []OpeningType    `json:"opening_types"    validate:"dive,omitempty,validate_opening_type"`
	CompanyDomains  []string         `json:"company_domains"  validate:"dive,omitempty,validate_domain,max=4"`
	ExperienceRange *ExperienceRange `json:"experience_range" validate:"omitempty"`
	SalaryRange     *SalaryRange     `json:"salary_range"     validate:"omitempty"`

	MinEducationLevel  *EducationLevel `json:"min_education_level"  validate:"omitempty,validate_education_level"`
	RemoteTimezones    []TimeZone      `json:"remote_timezones"     validate:"dive,omitempty,validate_timezone,max=4"`
	RemoteCountryCodes []CountryCode   `json:"remote_country_codes" validate:"dive,omitempty,validate_country_code,max=4"`
	PaginationKey      int64           `json:"pagination_key"`
	Limit              int64           `json:"limit"                validate:"min=1,max=100"`
}

type HubOpening struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company"`
	CompanyDomain          string `json:"company_domain"`
	CompanyName            string `json:"company_name"`
	JobTitle               string `json:"job_title"`
	JD                     string `json:"jd"`
	PaginationKey          int64  `json:"pagination_key"`
}
