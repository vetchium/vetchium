package db

import (
	"time"

	"github.com/google/uuid"
)

type EmployerTokenType string

const (
	// These are the session tokens for the employee. LTS refers to
	// Long Term Session and can be valid for a long time.
	EmployerSessionToken EmployerTokenType = "EMPLOYER_SESSION"
	EmployerLTSToken     EmployerTokenType = "EMPLOYER_LTS"

	// This is sent as a response to the signin request and should be used
	// in the tfa request, to get one of the session tokens.
	EmployerTFAToken EmployerTokenType = "EMPLOYER_TFA_TOKEN"
)

type HubTokenType string

const (
	// These are the session tokens for the HubUser. LTS refers to
	// Long Term Session and can be valid for a long time.
	HubUserSessionToken HubTokenType = "HUB_USER_SESSION"
	HubUserLTSToken     HubTokenType = "HUB_USER_LTS"

	// This is sent as a response to the Login request and should be used
	// in the tfa request, to get one of the session tokens.
	HubUserTFAToken HubTokenType = "HUB_USER_TFA_TOKEN"

	// This is sent as a response to the Reset Password request.
	HubUserResetPasswordToken HubTokenType = "HUB_USER_RESET_PASSWORD_TOKEN"
)

type HubTokenReq struct {
	Token            string
	TokenType        HubTokenType
	ValidityDuration time.Duration
	HubUserID        uuid.UUID
}

type EmployerTokenReq struct {
	Token            string
	TokenType        EmployerTokenType
	ValidityDuration time.Duration
	OrgUserID        uuid.UUID
}

type OrgUserInviteReq struct {
	Token            string
	ValidityDuration time.Duration
}

type HubUserInviteReq struct {
	Token              string
	ValidityDuration   time.Duration
	ReferringHubUserID uuid.UUID
}
