package vetchi

type HubUserState string

const (
	ActiveHubUserState   HubUserState = "ACTIVE_HUB_USER"
	DeletedHubUserState  HubUserState = "DELETED_HUB_USER"
	DisabledHubUserState HubUserState = "DISABLED_HUB_USER"
)

type LoginRequest struct {
	Email    EmailAddress `json:"email"    validate:"required,email"`
	Password Password     `json:"password" validate:"required,password"`
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

type InviteUserRequest struct {
	Email EmailAddress `json:"email" validate:"required,email"`
}

type GetMyHandleResponse struct {
	Handle string `json:"handle"`
}

type ChangePasswordRequest struct {
	OldPassword Password `json:"old_password" validate:"required,password"`
	NewPassword Password `json:"new_password" validate:"required,password"`
}

type ForgotPasswordRequest struct {
	Email EmailAddress `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Token    string   `json:"token"    validate:"required"`
	Password Password `json:"password" validate:"required,password"`
}
