package hub

import (
	"time"

	"github.com/vetchium/vetchium/typespec/common"
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

type ColleagueConnectionState string

const (
	CanSendRequestCCState         ColleagueConnectionState = "CAN_SEND_REQUEST"
	CannotSendRequestCCState      ColleagueConnectionState = "CANNOT_SEND_REQUEST"
	RequestSentPendingCCState     ColleagueConnectionState = "REQUEST_SENT_PENDING"
	RequestReceivedPendingCCState ColleagueConnectionState = "REQUEST_RECEIVED_PENDING"
	ConnectedCCState              ColleagueConnectionState = "CONNECTED"
	RejectedByMeCCState           ColleagueConnectionState = "REJECTED_BY_ME"
	RejectedByThemCCState         ColleagueConnectionState = "REJECTED_BY_THEM"
	UnlinkedByMeCCState           ColleagueConnectionState = "UNLINKED_BY_ME"
	UnlinkedByThemCCState         ColleagueConnectionState = "UNLINKED_BY_THEM"
)

type Bio struct {
	Handle                   string                   `json:"handle"`
	FullName                 string                   `json:"full_name"`
	ShortBio                 string                   `json:"short_bio"`
	LongBio                  string                   `json:"long_bio"`
	VerifiedMailDomains      []string                 `json:"verified_mail_domains"`
	ColleagueConnectionState ColleagueConnectionState `json:"colleague_connection_state"`
}

type UpdateBioRequest struct {
	FullName *string `json:"full_name" validate:"required,min=1,max=64"`
	ShortBio *string `json:"short_bio" validate:"required,min=1,max=64"`
	LongBio  *string `json:"long_bio"  validate:"required,min=1,max=1024"`
}

type UploadProfilePictureRequest struct {
	Image []byte `json:"image" validate:"required"`
}
