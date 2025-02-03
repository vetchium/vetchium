package employer

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type OpeningID string

type OpeningInfo struct {
	ID              string              `json:"id"               db:"id"`
	Title           string              `json:"title"            db:"title"`
	Positions       int                 `json:"positions"        db:"positions"`
	FilledPositions int                 `json:"filled_positions" db:"filled_positions"`
	Recruiter       OrgUserShort        `json:"recruiter"        db:"recruiter"`
	HiringManager   OrgUserShort        `json:"hiring_manager"   db:"hiring_manager"`
	CostCenterName  CostCenterName      `json:"cost_center_name" db:"cost_center_name"`
	OpeningType     common.OpeningType  `json:"opening_type"     db:"opening_type"`
	State           common.OpeningState `json:"state"            db:"state"`
	CreatedAt       time.Time           `json:"created_at"       db:"created_at"`
	LastUpdatedAt   time.Time           `json:"last_updated_at"  db:"last_updated_at"`
}

type Opening struct {
	ID                 string               `json:"id"`
	Title              string               `json:"title"`
	Positions          int                  `json:"positions"`
	FilledPositions    int                  `json:"filled_positions"`
	JD                 string               `json:"jd"`
	Recruiter          OrgUserShort         `json:"recruiter"`
	HiringManager      OrgUserShort         `json:"hiring_manager"`
	HiringTeam         []OrgUserShort       `json:"hiring_team,omitempty"`
	CostCenterName     CostCenterName       `json:"cost_center_name"`
	LocationTitles     []string             `json:"location_titles,omitempty"`
	RemoteCountryCodes []common.CountryCode `json:"remote_country_codes,omitempty"`
	RemoteTimezones    []common.TimeZone    `json:"remote_timezones,omitempty"`
	OpeningType        common.OpeningType   `json:"opening_type"`
	YoeMin             int                  `json:"yoe_min"`
	YoeMax             int                  `json:"yoe_max"`
	State              common.OpeningState  `json:"state"`
	CreatedAt          time.Time            `json:"created_at"`
	LastUpdatedAt      time.Time            `json:"last_updated_at"`

	// Optional fields
	EmployerNotes     *string               `json:"employer_notes,omitempty"`
	MinEducationLevel common.EducationLevel `json:"min_education_level"`
	Salary            *common.Salary        `json:"salary,omitempty"`
	Tags              []common.OpeningTag   `json:"tags,omitempty"`
}

type CreateOpeningRequest struct {
	Title          string                `json:"title"                     validate:"required,min=3,max=32"`
	Positions      int                   `json:"positions"                 validate:"required,min=1,max=20"`
	JD             string                `json:"jd"                        validate:"required,min=10,max=1024"`
	Recruiter      common.EmailAddress   `json:"recruiter"                 validate:"required"`
	HiringManager  common.EmailAddress   `json:"hiring_manager"            validate:"required"`
	HiringTeam     []common.EmailAddress `json:"hiring_team,omitempty"     validate:"omitempty,max=10"`
	CostCenterName CostCenterName        `json:"cost_center_name"          validate:"required"`
	LocationTitles []string              `json:"location_titles,omitempty" validate:"omitempty,max=10"`

	// TODO: Add validation for remote_country_codes and remote_timezones
	RemoteCountryCodes []common.CountryCode `json:"remote_country_codes,omitempty" validate:"omitempty,dive,validate_country_code,max=100"`
	RemoteTimezones    []common.TimeZone    `json:"remote_timezones,omitempty"     validate:"omitempty,max=200"`

	OpeningType common.OpeningType `json:"opening_type" validate:"required,validate_opening_type"`
	YoeMin      int                `json:"yoe_min"      validate:"min=0,max=100"`
	YoeMax      int                `json:"yoe_max"      validate:"min=1,max=100"`

	// Optional fields
	EmployerNotes     *string               `json:"employer_notes,omitempty" validate:"omitempty,max=1024"`
	MinEducationLevel common.EducationLevel `json:"min_education_level"      validate:"required,validate_education_level"`
	Salary            *common.Salary        `json:"salary,omitempty"         validate:"omitempty"`

	// Optional fields
	Tags    []common.OpeningTagID `json:"tags,omitempty"     validate:"omitempty,max=3"`
	NewTags []string              `json:"new_tags,omitempty" validate:"omitempty,max=3"`
}

type CreateOpeningResponse struct {
	OpeningID string `json:"opening_id"`
}

type GetOpeningRequest struct {
	ID string `json:"id" validate:"required"`
}

type FilterOpeningsRequest struct {
	State []common.OpeningState `json:"state,omitempty" validate:"omitempty,validate_opening_states"`

	FromDate *time.Time `json:"from_date,omitempty" validate:"omitempty,validate_opening_filter_start_date"`
	ToDate   *time.Time `json:"to_date,omitempty"   validate:"omitempty,validate_opening_filter_end_date"`

	PaginationKey string `json:"pagination_key,omitempty"`
	Limit         int    `json:"limit"                    validate:"max=40"`
}

func (filterOpeningsReq FilterOpeningsRequest) StatesAsStrings() []string {
	states := make([]string, len(filterOpeningsReq.State))
	for i, state := range filterOpeningsReq.State {
		// Already validated by validate_opening_state Vator
		states[i] = string(state)
	}
	return states
}

type ChangeOpeningStateRequest struct {
	OpeningID string              `json:"opening_id" validate:"required"`
	FromState common.OpeningState `json:"from_state" validate:"required"`
	ToState   common.OpeningState `json:"to_state"   validate:"required"`
}

type UpdateOpeningRequest struct {
	OpeningID string `json:"opening_id" validate:"required"`
	// TODO: Decide what fields are allowed to be updated
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
