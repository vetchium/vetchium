package vetchi

import (
	"time"
)

type AddWorkHistoryRequest struct {
	CompanyHandle string `json:"company_handle"     validate:"required"`
	JobTitle      string `json:"job_title"          validate:"required"`
	StartDate     string `json:"start_date"         validate:"required,date"`
	EndDate       string `json:"end_date,omitempty" validate:"date"`
}

type AddWorkHistoryResponse struct {
	WorkHistoryID string `json:"work_history_id"`
}

type Application struct {
	ApplicationID     string           `json:"application_id"`
	OpeningID         string           `json:"opening_id"`
	OpeningTitle      string           `json:"opening_title"`
	CompanyName       string           `json:"company_name"`
	CompanyHandle     string           `json:"company_handle"`
	CompanyLogo       string           `json:"company_logo,omitempty" validate:"url"`
	AppliedAt         time.Time        `json:"applied_at"`
	ApplicationState  ApplicationState `json:"application_state"`
	CandidacyID       string           `json:"candidacy_id,omitempty"`
	LastStateChangeAt time.Time        `json:"last_statechange_at"`
}

type ApplicationLabelColor string

const (
	Red     ApplicationLabelColor = "RED"
	Green   ApplicationLabelColor = "GREEN"
	Blue    ApplicationLabelColor = "BLUE"
	Magenta ApplicationLabelColor = "MAGENTA"
)

type ApplicationState string

const (
	ApplicationApplied     ApplicationState = "APPLICATION_APPLIED"
	ApplicationShortlisted ApplicationState = "APPLICATION_SHORTLISTED"
	ApplicationRejected    ApplicationState = "APPLICATION_REJECTED"
)

type ApplyToOpeningRequest struct {
	OpeningID       string   `json:"opening_id"                 validate:"required"`
	Resume          string   `json:"resume"                     validate:"required"`
	VouchersHandles []string `json:"vouchers_handles,omitempty"`
}

type ApplyToOpeningResponse struct {
	ApplicationID string `json:"application_id"`
}

type AutoBiography struct {
	NameEn                string   `json:"name_en"                            validate:"required,min=1,max=128"`
	AboutMe               string   `json:"about_me"                           validate:"required,min=1,max=1024"`
	Websites              []string `json:"websites,omitempty"                 validate:"omitempty,min=1,max=255,url"`
	NamesInOtherLanguages []struct {
		Language string `json:"language" validate:"required"`
		Name     string `json:"name" validate:"required"`
	} `json:"names_in_other_languages,omitempty"`
}

type CancelInterviewRequest struct {
	InterviewID                 string `json:"interview_id"                             validate:"required"`
	SendCancellationToCandidate bool   `json:"send_cancellation_to_candidate,omitempty"`
	CancellationBody            string `json:"cancellation_body,omitempty"              validate:"min=10,max=1000"`
}

type Candidacy struct {
	CandidacyID         string               `json:"candidacy_id"`
	CandidacyState      CandidacyState       `json:"candidacy_state"`
	Interviews          []MyInterview        `json:"interviews"`
	LastCompany         string               `json:"last_company"`
	LastPosition        string               `json:"last_position"`
	Name                string               `json:"name"`
	ShortlistedOpenings []ShortlistedOpening `json:"shortlisted_openings"`
}

type CandidacyState string

const (
	CandidacyShortlisted CandidacyState = "CANDIDACY_SHORTLISTED"
	CandidacyOffered     CandidacyState = "CANDIDACY_OFFERED"
	CandidacyRejected    CandidacyState = "CANDIDACY_REJECTED"
	CandidacyAccepted    CandidacyState = "CANDIDACY_ACCEPTED"
	CandidacyWithdrawn   CandidacyState = "CANDIDACY_WITHDRAWN"
	CandidacyCompleted   CandidacyState = "CANDIDACY_COMPLETED"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type ClientId struct {
	ClientID string `json:"client_id" validate:"required,min=3,max=255"`
}

type CreateInterviewRequest struct {
	CandidacyID              string    `json:"candidacy_id"                         validate:"required"`
	Interviewers             []OrgUser `json:"interviewers"                         validate:"required"`
	StartTime                time.Time `json:"start_time"                           validate:"required"`
	EndTime                  time.Time `json:"end_time"                             validate:"required"`
	SendAppointment          bool      `json:"send_appointment,omitempty"`
	InterviewAppointmentBody string    `json:"interview_appointment_body,omitempty" validate:"min=10,max=1000"`
}

type CreateInterviewResponse struct {
	InterviewID string `json:"interview_id"`
}

type EmailAddress string

type EmployerSignInRequest struct {
	ClientID string       `json:"client_id" validate:"required,client_id"`
	Email    EmailAddress `json:"email"     validate:"required,email"`
	Password Password     `json:"password"  validate:"required,password"`
}

type EmployerSignInResponse struct {
	Token string `json:"token"`
}

type EmployerTFARequest struct {
	TFACode    string `json:"tfa_code"              validate:"required"`
	TFAToken   string `json:"tfa_token"             validate:"required"`
	RememberMe bool   `json:"remember_me,omitempty"`
}

type EmployerTFAResponse struct {
	SessionToken string `json:"session_token"`
}

type EvaluationReport struct {
	Positives string `json:"positives" validate:"required,min=10,max=2048"`
	Negatives string `json:"negatives" validate:"required,min=10,max=2048"`
	Summary   string `json:"summary"   validate:"required,min=10,max=2048"`
}

type EvaluationResult string

const (
	StrongYes EvaluationResult = "STRONG_YES"
	Yes       EvaluationResult = "YES"
	No        EvaluationResult = "NO"
	StrongNo  EvaluationResult = "STRONG_NO"
)

type EvaluationState string

const (
	EvaluationPending   EvaluationState = "EVALUATION_PENDING"
	EvaluationCompleted EvaluationState = "EVALUATION_COMPLETED"
)

type FilterApplicationsRequest struct {
	ApplicationState ApplicationState        `json:"application_state,omitempty"`
	OpeningID        string                  `json:"opening_id"                  validate:"required"`
	SearchPrefix     string                  `json:"search_prefix,omitempty"`
	ColorFilters     []ApplicationLabelColor `json:"color_filters,omitempty"`
	Limit            int                     `json:"limit,omitempty"             validate:"min=1,max=100"`
	Offset           string                  `json:"offset,omitempty"`
}

type FilterCompaniesRequest struct {
	LanguageID string `json:"language_id"      validate:"required"`
	Prefix     string `json:"prefix"           validate:"required"`
	Offset     string `json:"offset,omitempty"`
	Limit      int    `json:"limit,omitempty"  validate:"min=1,max=10"`
}

type FilterEmployeesRequest struct {
	Prefix string `json:"prefix" validate:"required,min=3,max=255"`
}

type FilterHiringManagersRequest struct {
	Prefix string `json:"prefix" validate:"required"`
}

type FilterRecruitersRequest struct {
	Prefix string `json:"prefix" validate:"required"`
}

type FilteredApplication struct {
	ApplicationID string                  `json:"application_id"`
	ApplicantName string                  `json:"applicant_name"`
	LastCompany   string                  `json:"last_company"`
	LastPosition  string                  `json:"last_position"`
	VetchiHandle  string                  `json:"vetchi_handle"`
	ColorFilters  []ApplicationLabelColor `json:"color_filters,omitempty"`
	ResumeURL     string                  `json:"resume_url,omitempty"    validate:"url"`
	Vouches       []Vouch                 `json:"vouches,omitempty"`
	ReferredBy    ReferredBy              `json:"referred_by,omitempty"`
}

type FilteredCompany struct {
	CompanyHandle string `json:"company_handle"`
	CompanyName   string `json:"company_name"`
}

type FilteredEmployee struct {
	Name         string `json:"name"`
	Email        string `json:"email"        validate:"email"`
	VetchiHandle string `json:"vetchiHandle"`
}

type FilteredHiringManagers struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
}

type FilteredOpenings struct {
	DepartmentID       string       `json:"department_id"`
	DepartmentName     string       `json:"department_name"`
	HiringManagerEmail string       `json:"hiring_manager_email" validate:"email"`
	HiringManagerName  string       `json:"hiring_manager_name"`
	ID                 string       `json:"id"`
	RecruiterEmail     string       `json:"recruiter_email"      validate:"email"`
	RecruiterName      string       `json:"recruiter_name"`
	Status             OpeningState `json:"status"`
	Title              string       `json:"title"`
}

type FilteredRecruiters struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type GetInterviewDetailsRequest struct {
	InterviewID string `json:"interview_id" validate:"required"`
}

type GetMyApplicationsRequest struct {
	StartDate string `json:"start_date,omitempty" validate:"date"`
	EndDate   string `json:"end_date,omitempty"   validate:"date"`
}

type GetMyCandidaciesRequest struct {
	StartDate string `json:"start_date,omitempty" validate:"date"`
	EndDate   string `json:"end_date,omitempty"   validate:"date"`
}

type GetOnboardStatusRequest struct {
	ClientID string `json:"client_id" validate:"required,client_id"`
}

type GetOnboardStatusResponse struct {
	Status OnboardStatus `json:"status"`
}

type GetUserRequest struct {
	Limit  int    `json:"limit,omitempty"  validate:"min=1,max=100"`
	Offset string `json:"offset,omitempty" validate:"email"`
}

type GetWorkHistoryRequest struct {
	LanguageID string `json:"language_id" validate:"required"`
}

type HubTFAuthRequest struct {
	TFACode    string `json:"tfa_code"              validate:"required"`
	RememberMe bool   `json:"remember_me,omitempty"`
}

type HubTFAuthResponse struct {
	SessionToken string `json:"session_token"`
}

type Interview struct {
	InterviewID         string         `json:"interview_id"`
	InterviewState      InterviewState `json:"interview_status"`
	InterviewStartTime  time.Time      `json:"interview_start_time"`
	InterviewEndTime    time.Time      `json:"interview_end_time"`
	InterviewersNames   []string       `json:"interviewers_names"`
	FeedbackToCandidate string         `json:"feedback_to_candidate,omitempty"`
}

type InterviewDetails struct {
	InterviewID             string               `json:"interview_id"`
	InterviewState          InterviewState       `json:"interview_state"`
	Openings                []ShortlistedOpening `json:"openings"`
	CandidacyID             string               `json:"candidacy_id"`
	StartTime               time.Time            `json:"start_time"`
	EndTime                 time.Time            `json:"end_time"`
	CandidateName           string               `json:"candidate_name"`
	CandidateCurrentCompany string               `json:"candidate_current_company"`
	Interviewers            []OrgUser            `json:"interviewers"`
	EvaluationState         EvaluationState      `json:"evaluation_state"`
	EvaluationReport        EvaluationReport     `json:"evaluation_report"`
	EvaluationResult        EvaluationResult     `json:"evaluation_result"`
	FeedbackToCandidate     string               `json:"feedback_to_candidate,omitempty"`
}

type InterviewState string

const (
	InterviewScheduled InterviewState = "INTERVIEW_SCHEDULED"
	InterviewCompleted InterviewState = "INTERVIEW_COMPLETED"
	InterviewCancelled InterviewState = "INTERVIEW_CANCELLED"
	InterviewClosed    InterviewState = "INTERVIEW_CLOSED"
)

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	TFAToken  string    `json:"tfa_token"`
	ValidTill time.Time `json:"valid_till"`
}

type MyCandidacy struct {
	CandidacyID             string              `json:"candidacy_id"`
	Openings                []OpeningPublicInfo `json:"openings"`
	CandidacyState          CandidacyState      `json:"candidacy_state"`
	CandidacyStateUpdatedAt time.Time           `json:"candidacy_state_updated_at"`
	Interviews              []MyInterview       `json:"interviews"`
	CreatedAt               time.Time           `json:"created_at"`
}

type MyInterview struct {
	InterviewID         string         `json:"interview_id"`
	InterviewState      InterviewState `json:"interview_status"`
	InterviewStartTime  time.Time      `json:"interview_start_time"`
	InterviewEndTime    time.Time      `json:"interview_end_time"`
	InterviewersNames   []string       `json:"interviewers_names"`
	FeedbackToCandidate string         `json:"feedback_to_candidate,omitempty"`
}

type OnboardStatus string

const (
	DomainNotVerified            OnboardStatus = "DOMAIN_NOT_VERIFIED"
	DomainVerifiedOnboardPending OnboardStatus = "DOMAIN_VERIFIED_ONBOARD_PENDING"
	DomainOnboarded              OnboardStatus = "DOMAIN_ONBOARDED"
)

type OpeningInfo struct {
	DepartmentID       string       `json:"department_id"`
	DepartmentName     string       `json:"department_name"`
	FilledCount        int          `json:"filled_count"`
	HiringManagerEmail string       `json:"hiring_manager_email" validate:"email"`
	HiringManagerName  string       `json:"hiring_manager_name"`
	ID                 string       `json:"id"`
	JobType            string       `json:"job_type"`
	RecruiterEmail     string       `json:"recruiter_email"      validate:"email"`
	RecruiterName      string       `json:"recruiter_name"`
	OpeningState       OpeningState `json:"opening_state"`
	Title              string       `json:"title"`
	UnfilledCount      int          `json:"unfilled_count"`
}

type OpeningPublicInfo struct {
	OpeningID     string `json:"opening_id"`
	Title         string `json:"title"`
	CompanyName   string `json:"company_name"`
	CompanyHandle string `json:"company_handle"`
	CompanyLogo   string `json:"company_logo,omitempty" validate:"url"`
	JD            string `json:"jd"`
}

type Password string

type ReferredBy struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"email"`
}

type RemoveWorkHistoryRequest struct {
	WorkHistoryID string `json:"work_history_id" validate:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"        validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type SetOnboardPasswordRequest struct {
	ClientID string   `json:"client_id" validate:"required,client_id"`
	Password Password `json:"password"  validate:"required,password"`
	Token    string   `json:"token"     validate:"required"`
}

type ShortlistedOpening struct {
	OpeningID          string `json:"opening_id"`
	Title              string `json:"title"`
	HiringManagerName  string `json:"hiring_manager_name"`
	HiringManagerEmail string `json:"hiring_manager_email"`
	RecruiterName      string `json:"recruiter_name"`
	RecruiterEmail     string `json:"recruiter_email"`
}

type UpdateInterviewFeedbackRequest struct {
	InterviewID         string           `json:"interview_id"                    validate:"required"`
	EvaluationReport    EvaluationReport `json:"evaluation_report"               validate:"required"`
	EvaluationResult    EvaluationResult `json:"evaluation_result"               validate:"required"`
	FeedbackToCandidate string           `json:"feedback_to_candidate,omitempty" validate:"min=10,max=1000"`
}

type UpdateOrgUserRolesRequest struct {
	Email string        `json:"email" validate:"required,email"`
	Roles []OrgUserRole `json:"roles" validate:"required"`
}

type UpdateWorkHistoryRequest struct {
	WorkHistoryID string `json:"work_history_id"    validate:"required"`
	CompanyHandle string `json:"company_handle"     validate:"required"`
	JobTitle      string `json:"job_title"          validate:"required"`
	StartDate     string `json:"start_date"         validate:"required,date"`
	EndDate       string `json:"end_date,omitempty" validate:"date"`
}

type UpdateWorkHistoryResponse struct {
	WorkHistoryID string `json:"work_history_id"`
}

type ValidationErrors struct {
	Errors []string `json:"errors"`
}

type VetchiHandle struct {
	Handle string `json:"handle" validate:"required,min=6,max=32,pattern=^[a-zA-Z0-9_-]+$"`
}

type Vouch struct {
	VoucherName            string     `json:"voucher_name"`
	VoucherVetchiHandle    string     `json:"voucher_vetchi_handle"`
	VoucherCurrentPosition string     `json:"voucher_current_position"`
	VoucherCurrentCompany  string     `json:"voucher_current_company"`
	VouchState             VouchState `json:"vouch_state"`
}

type VouchState string

const (
	VouchSought   VouchState = "VOUCH_SOUGHT"
	VouchVouched  VouchState = "VOUCH_VOUCHED"
	VouchRejected VouchState = "VOUCH_REJECTED"
)

type WorkHistory struct {
	WorkHistoryID string `json:"work_history_id"`
	CompanyName   string `json:"company_name"`
	JobTitle      string `json:"job_title"`
	StartDate     string `json:"start_date"         validate:"date"`
	EndDate       string `json:"end_date,omitempty" validate:"date"`
	Logo          string `json:"logo,omitempty"     validate:"url"`
}
