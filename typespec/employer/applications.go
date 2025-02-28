package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
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
	State          common.ApplicationState `json:"state"            validate:"validate_application_state"`
	SearchQuery    *string                 `json:"search_query"     validate:"omitempty,max=25"`
	ColorTagFilter *ApplicationColorTag    `json:"color_tag_filter" validate:"omitempty"`
	OpeningID      string                  `json:"opening_id"       validate:"required"`
	PaginationKey  *string                 `json:"pagination_key"   validate:"omitempty"`
	Limit          int64                   `json:"limit"            validate:"required,min=0,max=40"`
}

type Endorser struct {
	FullName              string   `json:"full_name"`
	ShortBio              string   `json:"short_bio"`
	Handle                string   `json:"handle"`
	CurrentCompanyDomains []string `json:"current_company_domains"`
}

type Application struct {
	ID                        string                  `json:"id"`
	CoverLetter               *string                 `json:"cover_letter,omitempty"`
	CreatedAt                 time.Time               `json:"created_at"`
	HubUserHandle             string                  `json:"hub_user_handle"`
	HubUserName               string                  `json:"hub_user_name"`
	HubUserLastEmployerDomain *string                 `json:"hub_user_last_employer_domain,omitempty"`
	State                     common.ApplicationState `json:"state"`
	ColorTag                  *ApplicationColorTag    `json:"color_tag,omitempty"`
	Endorsers                 []Endorser              `json:"endorsers"`
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

type GetResumeRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}
