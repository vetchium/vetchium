package libvetchi

import "time"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	TFAToken  string    `json:"tfa_token"`
	ValidTill time.Time `json:"valid_till"`
}
