package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type GetApplicationsRequest struct {
	State          common.ApplicationState     `json:"state"            validate:"validate_application_state"`
	SearchQuery    *string                     `json:"search_query"     validate:"omitempty,max=25"`
	ColorTagFilter *common.ApplicationColorTag `json:"color_tag_filter" validate:"omitempty"`
	OpeningID      string                      `json:"opening_id"       validate:"required"`
	PaginationKey  *string                     `json:"pagination_key"   validate:"omitempty"`
	Limit          int64                       `json:"limit"            validate:"required,min=0,max=40"`
}

type Application struct {
	ID                        string                  `json:"id"`
	CoverLetter               *string                 `json:"cover_letter,omitempty"`
	CreatedAt                 time.Time               `json:"created_at"`
	Filename                  string                  `json:"filename"`
	HubUserHandle             string                  `json:"hub_user_handle"`
	HubUserLastEmployerDomain *string                 `json:"hub_user_last_employer_domain,omitempty"`
	Resume                    string                  `json:"resume"`
	State                     common.ApplicationState `json:"state"`
}

type UpdateApplicationStateRequest struct {
	ID        string                  `json:"id"         validate:"required"`
	FromState common.ApplicationState `json:"from_state" validate:"required"`
	ToState   common.ApplicationState `json:"to_state"   validate:"required"`
}

type SetApplicationColorTagRequest struct {
	ApplicationID string                     `json:"application_id" validate:"required"`
	ColorTag      common.ApplicationColorTag `json:"color_tag"      validate:"required,validate_application_color_tag"`
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
