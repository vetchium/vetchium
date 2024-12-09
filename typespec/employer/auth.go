package employer

import "github.com/psankar/vetchi/typespec/common"

type GetOnboardStatusRequest struct {
	ClientID string `json:"client_id" validate:"required,client_id"`
}

type OnboardStatus string

const (
	DomainNotVerified            OnboardStatus = "DOMAIN_NOT_VERIFIED"
	DomainVerifiedOnboardPending OnboardStatus = "DOMAIN_VERIFIED_ONBOARD_PENDING"
	DomainOnboarded              OnboardStatus = "DOMAIN_ONBOARDED"
)

type GetOnboardStatusResponse struct {
	Status OnboardStatus `json:"status"`
}

type SetOnboardPasswordRequest struct {
	ClientID string          `json:"client_id" validate:"required,client_id"`
	Password common.Password `json:"password"  validate:"required,password"`
	Token    string          `json:"token"     validate:"required"`
}

type EmployerSignInRequest struct {
	ClientID string              `json:"client_id" validate:"required,client_id"`
	Email    common.EmailAddress `json:"email"     validate:"required,email"`
	Password common.Password     `json:"password"  validate:"required,password"`
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
