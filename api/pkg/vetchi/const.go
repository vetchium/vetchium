package vetchi

import "time"

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

const (
	// Sent in the email to the org users
	OrgUserInviteTokenLenBytes = 16
	// Sent in the invite email to new hub users
	HubUserInviteTokenLenBytes = 16
	// Remember to change the email templates too if the duration changes
	HubUserInviteTokenValidity time.Duration = 3 * 24 * time.Hour // 3 days
	// Triggered by the forgot password request and sent to the user's email
	PasswordResetTokenLenBytes = 16
	// Used for the session tokens, sent as response after /tfa
	SessionTokenLenBytes = 8
	// Sent as a response to the signin request
	// Used for the /employer/tfa request body
	TGTokenLenBytes = 32

	// Used for the email code that is sent to the user's email for tfa
	EmailTokenLenBytes = 2
	// Used for the code that is sent to the user's email for add official email
	AddOfficialEmailCodeLenBytes = 2

	ApplicationIDLenBytes = 16
	CandidacyIDLenBytes   = 16
	CommentIDLenBytes     = 12
	EndorsementIDLenBytes = 16
	InterviewIDLenBytes   = 16
	PostIDLenBytes        = 24
	ResumeIDLenBytes      = 12
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

const (
	PreferredLanguage = "en"
)

const (
	MaxApplicationsToScorePerBatch = 10
)

// Timer intervals for granger background jobs
const (
	PruneTokensInterval             = 1 * time.Minute
	CreateOnboardEmailsInterval     = 3 * time.Second
	PruneOfficialEmailCodesInterval = 5 * time.Minute
	MailSenderInterval              = 5 * time.Second
	ScoreApplicationsInterval       = 1 * time.Minute
)
