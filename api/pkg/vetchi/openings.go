package vetchi

import "time"

type OpeningType string

const (
	FullTime    OpeningType = "FULL_TIME"
	PartTime    OpeningType = "PART_TIME"
	Contract    OpeningType = "CONTRACT"
	Internship  OpeningType = "INTERNSHIP"
	Unspecified OpeningType = "UNSPECIFIED"
)

type EducationLevel string

const (
	Bachelor    EducationLevel = "BACHELOR"
	Master      EducationLevel = "MASTER"
	Doctorate   EducationLevel = "DOCTORATE"
	NotMatters  EducationLevel = "NOT_MATTERS"
	Unspecified EducationLevel = "UNSPECIFIED"
)

type OpeningState string

const (
	DraftOpening     OpeningState = "DRAFT_OPENING"
	ActiveOpening    OpeningState = "ACTIVE_OPENING"
	SuspendedOpening OpeningState = "SUSPENDED_OPENING"
	ClosedOpening    OpeningState = "CLOSED_OPENING"
)

type Salary struct {
	MinAmount decimal.Decimal `json:"min_amount" validate:"required,min=0"`
	MaxAmount decimal.Decimal `json:"max_amount" validate:"required,min=1"`
	Currency  Currency        `json:"currency"   validate:"required"`
}

type Opening struct {
	ID                   string          `json:"id"`
	Title                string          `json:"title"                            validate:"required,min=3,max=32"`
	Positions            int             `json:"positions"                        validate:"required,min=1,max=20"`
	FilledPositions      int             `json:"filled_positions"                 validate:"min=0,max=20"`
	JD                   string          `json:"jd"                               validate:"required,min=10,max=1024"`
	Recruiters           []string        `json:"recruiters"                       validate:"required,min=1,max=10"`
	HiringManager        EmailAddress    `json:"hiring_manager"                   validate:"required"`
	CostCenterName       CostCenterName  `json:"cost_center_name"                 validate:"required"`
	EmployerNotes        *string         `json:"employer_notes,omitempty"         validate:"omitempty,max=1024"`
	LocationTitles       []string        `json:"location_titles,omitempty"        validate:"omitempty,max=10"`
	RemoteCountryCodes   []CountryCode   `json:"remote_country_codes,omitempty"   validate:"omitempty,max=100"`
	RemoteTimezones      []TimeZone      `json:"remote_timezones,omitempty"       validate:"omitempty,max=200"`
	OpeningType          OpeningType     `json:"opening_type"                     validate:"required"`
	YoeMin               int             `json:"yoe_min"                          validate:"min=0,max=100"`
	YoeMax               int             `json:"yoe_max"                          validate:"min=1,max=100"`
	MinEducationLevel    *EducationLevel `json:"min_education_level,omitempty"`
	Salary               *Salary         `json:"salary,omitempty"`
	CurrentState         OpeningState    `json:"current_state"                    validate:"required"`
	ApprovalWaitingState *OpeningState   `json:"approval_waiting_state,omitempty"`
	HiringTeam           []string        `json:"hiring_team,omitempty"            validate:"omitempty,max=10"`
	CreatedAt            time.Time       `json:"created_at"`
	LastUpdatedAt        time.Time       `json:"last_updated_at"`
}

type CreateOpeningRequest struct {
	Title              string          `json:"title"                          validate:"required,min=3,max=32"`
	Positions          int             `json:"positions"                      validate:"required,min=1,max=20"`
	JD                 string          `json:"jd"                             validate:"required,min=10,max=1024"`
	Recruiters         []string        `json:"recruiters"                     validate:"required,min=1,max=10"`
	HiringManager      EmailAddress    `json:"hiring_manager"                 validate:"required"`
	CostCenterName     CostCenterName  `json:"cost_center_name"               validate:"required"`
	EmployerNotes      *string         `json:"employer_notes,omitempty"       validate:"omitempty,max=1024"`
	LocationTitles     []string        `json:"location_titles,omitempty"      validate:"omitempty,max=10"`
	RemoteCountryCodes []CountryCode   `json:"remote_country_codes,omitempty" validate:"omitempty,max=100"`
	RemoteTimezones    []TimeZone      `json:"remote_timezones,omitempty"     validate:"omitempty,max=200"`
	OpeningType        OpeningType     `json:"opening_type"                   validate:"required"`
	YoeMin             int             `json:"yoe_min"                        validate:"min=0,max=100"`
	YoeMax             int             `json:"yoe_max"                        validate:"min=1,max=100"`
	MinEducationLevel  *EducationLevel `json:"min_education_level,omitempty"`
	Salary             *Salary         `json:"salary,omitempty"`
}

type GetOpeningRequest struct {
	ID string `json:"id" validate:"required"`
}

type FilterOpeningsRequest struct {
	PaginationKey *string        `json:"pagination_key,omitempty"`
	State         []OpeningState `json:"state,omitempty"`
	Limit         *int           `json:"limit,omitempty"          validate:"omitempty,max=40"`
}

type UpdateOpeningRequest struct {
	ID                 string          `json:"id"                             validate:"required"`
	Title              string          `json:"title"                          validate:"required,min=3,max=32"`
	Positions          int             `json:"positions"                      validate:"required,min=1,max=20"`
	JD                 string          `json:"jd"                             validate:"required,min=10,max=1024"`
	Recruiters         []string        `json:"recruiters"                     validate:"required,min=1,max=10"`
	HiringManager      EmailAddress    `json:"hiring_manager"                 validate:"required"`
	CostCenterName     CostCenterName  `json:"cost_center_name"               validate:"required"`
	EmployerNotes      *string         `json:"employer_notes,omitempty"       validate:"omitempty,max=1024"`
	LocationTitles     []string        `json:"location_titles,omitempty"      validate:"omitempty,max=10"`
	RemoteCountryCodes []CountryCode   `json:"remote_country_codes,omitempty" validate:"omitempty,max=100"`
	RemoteTimezones    []TimeZone      `json:"remote_timezones,omitempty"     validate:"omitempty,max=200"`
	OpeningType        OpeningType     `json:"opening_type"                   validate:"required"`
	YoeMin             int             `json:"yoe_min"                        validate:"min=0,max=100"`
	YoeMax             int             `json:"yoe_max"                        validate:"min=1,max=100"`
	MinEducationLevel  *EducationLevel `json:"min_education_level,omitempty"`
	Salary             *Salary         `json:"salary,omitempty"`
}

type GetOpeningWatchersRequest struct {
	ID string `json:"id" validate:"required"`
}

type OpeningWatchers struct {
	ID     string         `json:"id"`
	Emails []EmailAddress `json:"emails,omitempty" validate:"omitempty,max=20"`
}

type AddWatchersRequest struct {
	ID     string         `json:"id"     validate:"required"`
	Emails []EmailAddress `json:"emails" validate:"required"`
}

type RemoveWatcherRequest struct {
	ID    string       `json:"id"    validate:"required"`
	Email EmailAddress `json:"email" validate:"required"`
}

type ApproveOpeningStateChangeRequest struct {
	ID string `json:"id" validate:"required"`
}

type RejectOpeningStateChangeRequest struct {
	ID string `json:"id" validate:"required"`
}
