package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type MyApplicationsRequest struct {
	State         common.ApplicationState `json:"state"          validate:"omitempty,validate_application_state"`
	PaginationKey *string                 `json:"pagination_key" validate:"omitempty"`
	Limit         int64                   `json:"limit"          validate:"required,min=0,max=40"`
}

type HubApplication struct {
	ApplicationID  string                  `json:"application_id"`
	State          common.ApplicationState `json:"state"`
	OpeningID      string                  `json:"opening_id"`
	OpeningTitle   string                  `json:"opening_title"`
	EmployerName   string                  `json:"employer_name"`
	EmployerDomain string                  `json:"employer_domain"`
	CreatedAt      time.Time               `json:"created_at"`
}

type WithdrawApplicationRequest struct {
	ApplicationID string `json:"application_id" validate:"required"`
}
