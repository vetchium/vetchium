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

type GetBioRequest struct {
	Handle string `json:"handle" validate:"required"`
}

type Bio struct {
	Handle              string   `json:"handle"`
	FullName            string   `json:"full_name"`
	ShortBio            string   `json:"short_bio"`
	LongBio             string   `json:"long_bio"`
	VerifiedMailDomains []string `json:"verified_mail_domains"`
	IsColleaguable      bool     `json:"is_colleaguable"`
	IsColleague         bool     `json:"is_colleague"`
}

type UpdateBioRequest struct {
	Handle   *string `json:"handle"    validate:"required,validate_handle,min=1,max=32"`
	FullName *string `json:"full_name" validate:"required,min=1,max=64"`
	ShortBio *string `json:"short_bio" validate:"required,min=1,max=64"`
	LongBio  *string `json:"long_bio"  validate:"required,min=1,max=1024"`
}

type UploadProfilePictureRequest struct {
	Image []byte `json:"image" validate:"required"`
}
