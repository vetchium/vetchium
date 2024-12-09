package hub

import "github.com/psankar/vetchi/typespec/common"

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
