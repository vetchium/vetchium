package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type AddOfficialEmailRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
}

type VerifyOfficialEmailRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
	Code  string              `json:"code"  validate:"required"`
}

type TriggerVerificationRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
}

type DeleteOfficialEmailRequest struct {
	Email common.EmailAddress `json:"email" validate:"required,email"`
}

type OfficialEmail struct {
	Email            common.EmailAddress `json:"email"`
	LastVerifiedAt   *time.Time          `json:"last_verified_at"`
	VerifyInProgress bool                `json:"verify_in_progress"`
}
