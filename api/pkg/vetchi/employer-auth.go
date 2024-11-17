package vetchi

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

type GetOnboardStatusRequest struct {
	ClientID string `json:"client_id" validate:"required,client_id"`
}

type GetOnboardStatusResponse struct {
	Status OnboardStatus `json:"status"`
}

type OnboardStatus string

const (
	DomainNotVerified            OnboardStatus = "DOMAIN_NOT_VERIFIED"
	DomainVerifiedOnboardPending OnboardStatus = "DOMAIN_VERIFIED_ONBOARD_PENDING"
	DomainOnboarded              OnboardStatus = "DOMAIN_ONBOARDED"
)

type SetOnboardPasswordRequest struct {
	ClientID string   `json:"client_id" validate:"required,client_id"`
	Password Password `json:"password"  validate:"required,password"`
	Token    string   `json:"token"     validate:"required"`
}
