package vetchi

import "time"

const (
	HubBaseURL      = "https://vetchi.org"
	EmployerBaseURL = "https://employer.vetchi.org"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

const (
	// Sent in the email to the org users
	InviteTokenLenBytes = 16

	// Triggered by the forgot password request and sent to the user's email
	PasswordResetTokenLenBytes = 16

	// Sent as a response to the signin request
	// Used for the /employer/tfa request body
	TGTokenLenBytes = 32

	// Used for the email code that is sent to the user's email for tfa
	EmailTokenLenBytes = 2

	// Used for the session tokens
	SessionTokenLenBytes = 8

	ApplicationIDLenBytes = 16

	CandidacyIDLenBytes = 16

	InterviewIDLenBytes = 16

	ResumeIDLenBytes = 12

	// Used for the code that is sent to the user's email for add official email
	AddOfficialEmailCodeLenBytes = 2
)

const (
	EmailFrom = "no-reply@vetchi.org"
)

const (
	// Duration for which an official email verification is considered valid
	VerificationValidityDuration = 90 * 24 * time.Hour // 90 days
)

const (
	// Profile picture constraints
	MaxProfilePictureSize    = 5 * 1024 * 1024 // 5MB
	MinProfilePictureDim     = 200
	MaxProfilePictureDim     = 2048
	ProfilePictureIDLenBytes = 16
)
