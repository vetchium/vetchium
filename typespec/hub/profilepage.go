package hub

import (
	"time"

	"github.com/psankar/vetchi/typespec/common"
)

type AddProfessionalEmailRequest struct {
	Email common.EmailAddress `json:"email"`
}

type VerifyProfessionalEmailRequest struct {
	Email common.EmailAddress `json:"email"`
}

type TriggerVerificationRequest struct {
	Email common.EmailAddress `json:"email"`
}

type DeleteProfessionalEmailRequest struct {
	Email common.EmailAddress `json:"email"`
}

type ProfessionalEmail struct {
	Email            common.EmailAddress `json:"email"`
	LastVerifiedAt   *time.Time          `json:"last_verified_at"`
	VerifyInProgress bool                `json:"verify_in_progress"`
}
