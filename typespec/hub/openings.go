package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type ExperienceRange struct {
	YoeMin int `json:"yoe_min" validate:"min=0,max=99"`
	YoeMax int `json:"yoe_max" validate:"min=1,max=100"`
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
	CountryCode common.CountryCode `json:"country_code"`
	Cities      []string           `json:"cities"       validate:"dive,omitempty"`

	OpeningTypes    []common.OpeningType `json:"opening_types"    validate:"dive,omitempty,validate_opening_type"`
	CompanyDomains  []string             `json:"company_domains"  validate:"dive,omitempty,validate_domain"`
	ExperienceRange *ExperienceRange     `json:"experience_range" validate:"omitempty"`
	SalaryRange     *common.Salary       `json:"salary_range"     validate:"omitempty"`

	MinEducationLevel *common.EducationLevel `json:"min_education_level" validate:"omitempty,validate_education_level"`
	PaginationKey     int64                  `json:"pagination_key"`
	Limit             int64                  `json:"limit"               validate:"min=0,max=100"`
}

type HubOpening struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company"`
	CompanyDomain          string `json:"company_domain"`
	CompanyName            string `json:"company_name"`
	JobTitle               string `json:"job_title"`
	JD                     string `json:"jd"`
	PaginationKey          int64  `json:"pagination_key"`
}

type GetHubOpeningDetailsRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
}

type GetHubOpeningDetailsResponse struct {
	CompanyDomain             string                 `json:"company_domain"`
	CompanyName               string                 `json:"company_name"`
	CreatedAt                 time.Time              `json:"created_at"`
	EducationLevel            *common.EducationLevel `json:"education_level"`
	ExperienceRange           *ExperienceRange       `json:"experience_range"`
	HiringManagerName         string                 `json:"hiring_manager_name"`
	HiringManagerVetchiHandle *string                `json:"hiring_manager_vetchi_handle"`
	JD                        string                 `json:"jd"`
	JobTitle                  string                 `json:"job_title"`
	OpeningIDWithinCompany    string                 `json:"opening_id_within_company"`
	OpeningType               common.OpeningType     `json:"opening_type"`
	PaginationKey             int64                  `json:"pagination_key"`
	RecruiterName             string                 `json:"recruiter_name"`
	Salary                    *common.Salary         `json:"salary"`
	State                     common.OpeningState    `json:"state"`
}

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
	Resume                 string `json:"resume"                    validate:"required"`
	Filename               string `json:"filename"                  validate:"required,max=256"`
	CoverLetter            string `json:"cover_letter"              validate:"omitempty,max=4096"`
}
