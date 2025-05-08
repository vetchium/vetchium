package hub

import "github.com/vetchium/vetchium/typespec/common"

type LoginRequest struct {
	Email    common.EmailAddress `json:"email"    validate:"required,email"`
	Password common.Password     `json:"password" validate:"required,password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type HubTFARequest struct {
	TFAToken   string `json:"tfa_token"             validate:"required"`
	TFACode    string `json:"tfa_code"              validate:"required"`
	RememberMe bool   `json:"remember_me,omitempty"`
}

type HubTFAResponse struct {
	SessionToken string `json:"session_token"`
}

type ChangePasswordRequest struct {
	OldPassword common.Password `json:"old_password" validate:"required,password"`
	NewPassword common.Password `json:"new_password" validate:"required,password"`
}

type ForgotPasswordRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string          `json:"token"    validate:"required"`
	Password common.Password `json:"password" validate:"required,password"`
}

type GetMyHandleResponse struct {
	Handle string `json:"handle"`
}
type HubUserState string

const (
	ActiveHubUserState HubUserState = "ACTIVE_HUB_USER"
)

type HubUserInviteRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
}

type HubUserTier string

const (
	FreeHubUserTier HubUserTier = "FREE_HUB_USER"
	PaidHubUserTier HubUserTier = "PAID_HUB_USER"
)

type OnboardHubUserRequest struct {
	Token               string             `json:"token"                 validate:"required"`
	FullName            string             `json:"full_name"             validate:"required"`
	ResidentCountryCode common.CountryCode `json:"resident_country_code" validate:"required,validate_country_code"`
	Password            common.Password    `json:"password"              validate:"required,password"`
	SelectedTier        HubUserTier        `json:"selected_tier"         validate:"required"`

	// TODO: Remove hard-coded language
	PreferredLanguage string `json:"preferred_language" validate:"required,eq=en"`

	// TODO: Make the lengths consistent across various APIs
	ShortBio string `json:"short_bio" validate:"required,max=64"`
	LongBio  string `json:"long_bio"  validate:"required,max=2048"`
}

type OnboardHubUserResponse struct {
	SessionToken    string `json:"session_token"    validate:"required"`
	GeneratedHandle string `json:"generated_handle" validate:"required"`
}

type CheckHandleAvailabilityRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}

type CheckHandleAvailabilityResponse struct {
	IsAvailable           bool     `json:"is_available"                     validate:"required"`
	SuggestedAlternatives []string `json:"suggested_alternatives,omitempty"`
}

type SetHandleRequest struct {
	Handle common.Handle `json:"handle" validate:"required,validate_handle"`
}
