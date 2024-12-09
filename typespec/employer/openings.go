package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type OpeningState string

const (
	DraftOpening     OpeningState = "DRAFT_OPENING_STATE"
	ActiveOpening    OpeningState = "ACTIVE_OPENING_STATE"
	SuspendedOpening OpeningState = "SUSPENDED_OPENING_STATE"
	ClosedOpening    OpeningState = "CLOSED_OPENING_STATE"
)

type Salary struct {
	MinAmount float64         `json:"min_amount" validate:"required,min=0"`
	MaxAmount float64         `json:"max_amount" validate:"required,min=1"`
	Currency  common.Currency `json:"currency"   validate:"required"`
}

type CreateOpeningRequest struct {
	Title              string                 `json:"title"                validate:"required,min=3,max=32"`
	Positions          int                    `json:"positions"            validate:"required,min=1,max=20"`
	JD                 string                 `json:"jd"                   validate:"required,min=10,max=1024"`
	RecruiterEmail     common.EmailAddress    `json:"recruiter_email"      validate:"required"`
	HiringManagerEmail common.EmailAddress    `json:"hiring_manager_email" validate:"required"`
	HiringTeamEmails   []common.EmailAddress  `json:"hiring_team_emails"   validate:"omitempty,max=10"`
	CostCenterName     CostCenterName         `json:"cost_center_name"     validate:"required"`
	LocationTitles     []string               `json:"location_titles"      validate:"required,min=1,dive,min=3,max=32"`
	RemoteCountryCodes []common.CountryCode   `json:"remote_country_codes" validate:"omitempty,dive,validate_country_code"`
	RemoteTimezones    []common.TimeZone      `json:"remote_timezones"     validate:"omitempty,dive,validate_timezone"`
	OpeningType        common.OpeningType     `json:"opening_type"         validate:"required,validate_opening_type"`
	YoeMin             int                    `json:"yoe_min"              validate:"min=0,max=99"`
	YoeMax             int                    `json:"yoe_max"              validate:"min=1,max=100"`
	MinEducationLevel  *common.EducationLevel `json:"min_education_level"  validate:"omitempty,validate_education_level"`
	Salary             *Salary                `json:"salary"               validate:"omitempty"`
}

type GetOpeningRequest struct {
	ID string `json:"id" validate:"required"`
}

type Opening struct {
	ID                 string                 `json:"id"`
	Title              string                 `json:"title"`
	Positions          int                    `json:"positions"`
	FilledPositions    int                    `json:"filled_positions"`
	JD                 string                 `json:"jd"`
	Recruiter          OrgUserShort           `json:"recruiter"`
	HiringManager      OrgUserShort           `json:"hiring_manager"`
	HiringTeam         []OrgUserShort         `json:"hiring_team,omitempty"`
	CostCenterName     CostCenterName         `json:"cost_center_name"`
	LocationTitles     []string               `json:"location_titles,omitempty"`
	RemoteCountryCodes []common.CountryCode   `json:"remote_country_codes,omitempty"`
	RemoteTimezones    []common.TimeZone      `json:"remote_timezones,omitempty"`
	OpeningType        common.OpeningType     `json:"opening_type"`
	YoeMin             int                    `json:"yoe_min"`
	YoeMax             int                    `json:"yoe_max"`
	State              OpeningState           `json:"state"`
	CreatedAt          time.Time              `json:"created_at"`
	LastUpdatedAt      time.Time              `json:"last_updated_at"`
	EmployerNotes      *string                `json:"employer_notes,omitempty"`
	MinEducationLevel  *common.EducationLevel `json:"min_education_level,omitempty"`
	Salary             *Salary                `json:"salary,omitempty"`
}

type OpeningInfo struct {
	ID              string       `json:"id"`
	Title           string       `json:"title"`
	Positions       int          `json:"positions"`
	FilledPositions int          `json:"filled_positions"`
	Recruiter       OrgUserShort `json:"recruiter"`
	HiringManager   OrgUserShort `json:"hiring_manager"`
	CostCenterName  string       `json:"cost_center_name"`
	OpeningType     string       `json:"opening_type"`
	State           OpeningState `json:"state"`
	CreatedAt       time.Time    `json:"created_at"`
	LastUpdatedAt   time.Time    `json:"last_updated_at"`
}

type FilterOpeningsRequest struct {
	State         []OpeningState `json:"state"           validate:"required,validate_opening_state"`
	FromDate      *time.Time     `json:"from_date"       validate:"omitempty"`
	ToDate        *time.Time     `json:"to_date"         validate:"omitempty"`
	PaginationKey string         `json:"pagination_key"`
	Limit         int            `json:"limit,omitempty" validate:"omitempty,max=100"`
}

type UpdateOpeningRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetOpeningWatchersRequest struct {
	OpeningID string `json:"opening_id" validate:"required"`
}

type AddOpeningWatchersRequest struct {
	OpeningID string                `json:"opening_id" validate:"required"`
	Emails    []common.EmailAddress `json:"emails"     validate:"required"`
}

type RemoveOpeningWatcherRequest struct {
	OpeningID string              `json:"opening_id" validate:"required"`
	Email     common.EmailAddress `json:"email"      validate:"required"`
}

type ChangeOpeningStateRequest struct {
	OpeningID string       `json:"opening_id" validate:"required"`
	FromState OpeningState `json:"from_state" validate:"required"`
	ToState   OpeningState `json:"to_state"   validate:"required"`
}
