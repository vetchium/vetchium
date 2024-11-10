package vetchi

import "time"

type OpeningType string

const (
	FullTimeOpening    OpeningType = "FULL_TIME_OPENING"
	PartTimeOpening    OpeningType = "PART_TIME_OPENING"
	ContractOpening    OpeningType = "CONTRACT_OPENING"
	InternshipOpening  OpeningType = "INTERNSHIP_OPENING"
	UnspecifiedOpening OpeningType = "UNSPECIFIED_OPENING"
)

type EducationLevel string

const (
	BachelorEducation    EducationLevel = "BACHELOR_EDUCATION"
	MasterEducation      EducationLevel = "MASTER_EDUCATION"
	DoctorateEducation   EducationLevel = "DOCTORATE_EDUCATION"
	NotMattersEducation  EducationLevel = "NOT_MATTERS_EDUCATION"
	UnspecifiedEducation EducationLevel = "UNSPECIFIED_EDUCATION"
)

type OpeningState string

const (
	DraftOpening     OpeningState = "DRAFT_OPENING"
	ActiveOpening    OpeningState = "ACTIVE_OPENING"
	SuspendedOpening OpeningState = "SUSPENDED_OPENING"
	ClosedOpening    OpeningState = "CLOSED_OPENING"
)

type Salary struct {
	MinAmount float64  `json:"min_amount" validate:"required,min=0"`
	MaxAmount float64  `json:"max_amount" validate:"required,min=1"`
	Currency  Currency `json:"currency"   validate:"required"`
}

type Opening struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	Positions          int            `json:"positions"`
	FilledPositions    int            `json:"filled_positions"`
	JD                 string         `json:"jd"`
	Recruiter          OrgUserShort   `json:"recruiter"`
	HiringManager      OrgUserShort   `json:"hiring_manager"`
	HiringTeam         []OrgUserShort `json:"hiring_team,omitempty"`
	CostCenterName     CostCenterName `json:"cost_center_name"`
	LocationTitles     []string       `json:"location_titles,omitempty"`
	RemoteCountryCodes []CountryCode  `json:"remote_country_codes,omitempty"`
	RemoteTimezones    []TimeZone     `json:"remote_timezones,omitempty"`
	OpeningType        OpeningType    `json:"opening_type"`
	YoeMin             int            `json:"yoe_min"`
	YoeMax             int            `json:"yoe_max"`
	CurrentState       OpeningState   `json:"current_state"`
	CreatedAt          time.Time      `json:"created_at"`
	LastUpdatedAt      time.Time      `json:"last_updated_at"`

	// Optional fields
	ApprovalWaitingState *OpeningState   `json:"approval_waiting_state,omitempty"`
	EmployerNotes        *string         `json:"employer_notes,omitempty"`
	MinEducationLevel    *EducationLevel `json:"min_education_level,omitempty"`
	Salary               *Salary         `json:"salary,omitempty"`
}

type CreateOpeningRequest struct {
	Title              string         `json:"title"                          validate:"required,min=3,max=32"`
	Positions          int            `json:"positions"                      validate:"required,min=1,max=20"`
	JD                 string         `json:"jd"                             validate:"required,min=10,max=1024"`
	Recruiter          EmailAddress   `json:"recruiter"                      validate:"required"`
	HiringManager      EmailAddress   `json:"hiring_manager"                 validate:"required"`
	HiringTeam         []EmailAddress `json:"hiring_team,omitempty"          validate:"omitempty,max=10"`
	CostCenterName     CostCenterName `json:"cost_center_name"               validate:"required"`
	LocationTitles     []string       `json:"location_titles,omitempty"      validate:"omitempty,max=10"`
	RemoteCountryCodes []CountryCode  `json:"remote_country_codes,omitempty" validate:"omitempty,max=100"`
	RemoteTimezones    []TimeZone     `json:"remote_timezones,omitempty"     validate:"omitempty,max=200"`
	OpeningType        OpeningType    `json:"opening_type"                   validate:"required"`
	YoeMin             int            `json:"yoe_min"                        validate:"min=0,max=100"`
	YoeMax             int            `json:"yoe_max"                        validate:"min=1,max=100"`

	// Optional fields
	EmployerNotes     *string         `json:"employer_notes,omitempty"      validate:"omitempty,max=1024"`
	MinEducationLevel *EducationLevel `json:"min_education_level,omitempty" validate:"omitempty"`
	Salary            *Salary         `json:"salary,omitempty"              validate:"omitempty"`
}

type GetOpeningRequest struct {
	ID string `json:"id" validate:"required"`
}

type FilterOpeningsRequest struct {
	State []OpeningState `json:"state,omitempty"`

	PaginationKey string `json:"pagination_key,omitempty"`
	Limit         int    `json:"limit,omitempty"          validate:"omitempty,max=40"`
}

type UpdateOpeningRequest struct {
	ID                string         `json:"id"                            validate:"required"`
	Title             string         `json:"title"                         validate:"required,min=3,max=32"`
	Positions         int            `json:"positions"                     validate:"required,min=1,max=20"`
	JD                string         `json:"jd"                            validate:"required,min=10,max=1024"`
	Recruiter         EmailAddress   `json:"recruiter"                     validate:"required"`
	HiringManager     EmailAddress   `json:"hiring_manager"                validate:"required"`
	CostCenterName    CostCenterName `json:"cost_center_name"              validate:"required"`
	OpeningType       OpeningType    `json:"opening_type"                  validate:"required"`
	EmployerNotes     string         `json:"employer_notes,omitempty"      validate:"omitempty,max=1024"`
	YoeMin            int            `json:"yoe_min"                       validate:"min=0,max=100"`
	YoeMax            int            `json:"yoe_max"                       validate:"min=1,max=100"`
	MinEducationLevel EducationLevel `json:"min_education_level,omitempty"`
	Salary            Salary         `json:"salary,omitempty"`
}

type GetOpeningWatchersRequest struct {
	ID string `json:"id" validate:"required"`
}

type OpeningWatchers struct {
	ID     string         `json:"id"`
	Emails []EmailAddress `json:"emails,omitempty" validate:"omitempty,max=20"`
}

type AddOpeningWatchersRequest struct {
	ID     string         `json:"id"     validate:"required"`
	Emails []EmailAddress `json:"emails" validate:"required"`
}

type RemoveOpeningWatcherRequest struct {
	ID    string       `json:"id"    validate:"required"`
	Email EmailAddress `json:"email" validate:"required"`
}

type ApproveOpeningStateChangeRequest struct {
	ID string `json:"id" validate:"required"`
}

type RejectOpeningStateChangeRequest struct {
	ID string `json:"id" validate:"required"`
}
