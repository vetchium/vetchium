package vetchi

import (
	"time"
)

type ApplyForOpeningRequest struct {
	OpeningIDWithinCompany string `json:"opening_id_within_company" validate:"required"`
	CompanyDomain          string `json:"company_domain"            validate:"required"`
	Resume                 string `json:"resume"                    validate:"required"`
	Filename               string `json:"filename"                  validate:"required,max=256"`
	CoverLetter            string `json:"cover_letter"              validate:"omitempty,max=4096"`
}

type ApplicationState string

const (
	AppliedAppState ApplicationState = "APPLIED"

	// TODO: Remember to Reject all open applications when an Opening is closed
	RejectedAppState ApplicationState = "REJECTED"

	ShortlistedAppState ApplicationState = "SHORTLISTED"
	WithdrawnAppState   ApplicationState = "WITHDRAWN"
	ExpiredAppState     ApplicationState = "EXPIRED"
)

type ApplicationColorTag string

const (
	// Any change here should reflect in the IsValid() method too
	GreenApplicationColorTag  ApplicationColorTag = "GREEN"
	YellowApplicationColorTag ApplicationColorTag = "YELLOW"
	RedApplicationColorTag    ApplicationColorTag = "RED"
)

func (c ApplicationColorTag) IsValid() bool {
	return c == GreenApplicationColorTag ||
		c == YellowApplicationColorTag ||
		c == RedApplicationColorTag
}

type GetApplicationsRequest struct {
	State          ApplicationState     `json:"state"            validate:"required"`
	SearchQuery    *string              `json:"search_query"     validate:"omitempty,max=25"`
	ColorTagFilter *ApplicationColorTag `json:"color_tag_filter" validate:"omitempty"`
	OpeningID      string               `json:"opening_id"       validate:"required"`
	PaginationKey  *string              `json:"pagination_key"   validate:"omitempty"`
	Limit          int64                `json:"limit"            validate:"required,min=0,max=40"`
}

type Application struct {
	ID                        string           `json:"id"`
	CoverLetter               *string          `json:"cover_letter,omitempty"`
	CreatedAt                 time.Time        `json:"created_at"`
	Filename                  string           `json:"filename"`
	HubUserHandle             string           `json:"hub_user_handle"`
	HubUserLastEmployerDomain *string          `json:"hub_user_last_employer_domain,omitempty"`
	Resume                    string           `json:"resume"`
	State                     ApplicationState `json:"state"`
}

type UpdateApplicationStateRequest struct {
	ID        string           `json:"id"         validate:"required"`
	FromState ApplicationState `json:"from_state" validate:"required"`
	ToState   ApplicationState `json:"to_state"   validate:"required"`
}

type SetApplicationColorTagRequest struct {
	ApplicationID string              `json:"application_id" validate:"required"`
	ColorTag      ApplicationColorTag `json:"color_tag"      validate:"required,validate_application_color_tag"`
}

type RemoveApplicationColorTagRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}

type ShortlistApplicationRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}

type RejectApplicationRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}

type MyApplicationsRequest struct {
	State         ApplicationState `json:"state"          validate:"required"`
	PaginationKey *string          `json:"pagination_key" validate:"omitempty"`
	Limit         int64            `json:"limit"          validate:"required,min=0,max=40"`
}

type HubApplication struct {
	ApplicationID  string           `json:"application_id"`
	State          ApplicationState `json:"state"`
	OpeningID      string           `json:"opening_id"`
	OpeningTitle   string           `json:"opening_title"`
	EmployerName   string           `json:"employer_name"`
	EmployerDomain string           `json:"employer_domain"`
	CreatedAt      time.Time        `json:"created_at"`
}

type WithdrawApplicationRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}
