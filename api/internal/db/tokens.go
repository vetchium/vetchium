package db

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	// These are the session tokens for the employee. LTS refers to
	// Long Term Session and can be valid for a long time.
	EmployerSessionToken TokenType = "EMPLOYER_SESSION"
	EmployerLTSToken     TokenType = "EMPLOYER_LTS"

	// This is sent as a response to the signin request and should be used
	// in the tfa request, to get one of the session tokens.
	EmployerTFAToken TokenType = "EMPLOYER_TFA_TOKEN"

	// This is emailed to the OrgUser after a sucessful signin request and
	// should be used in the tfa request as part of the body, to get one
	// of the session tokens.
	EmployerTFACode TokenType = "EMPLOYER_TFA_CODE"

	// This is emailed to a potential Employee, on behalf of their Employer,
	// to invite them to Vetchi so that they could become an OrgUser
	EmployerInviteToken TokenType = "EMPLOYER_INVITE"
)

const (
	// These are the session tokens for the HubUser. LTS refers to
	// Long Term Session and can be valid for a long time.
	HubSessionToken TokenType = "HUB_SESSION"
	HubLTSToken     TokenType = "HUB_LTS"

	// This is sent as a response to the Login request and should be used
	// in the tfa request, to get one of the session tokens.
	HubUserTFAToken TokenType = "HUB_USER_TFA_TOKEN"

	// This is emailed to the HubUser after a sucessful signin request and
	// should be used in the tfa request as part of the body, to get one
	// of the session tokens.
	HubUserTFACode TokenType = "HUB_USER_TFA_CODE"

	// This is emailed to a potential HubUser to invite them to the Hub.
	HubUserInviteToken TokenType = "HUB_USER_INVITE"
)

type TokenReq struct {
	Token            string
	TokenType        TokenType
	ValidityDuration time.Duration
	OrgUserID        uuid.UUID
	HubUserID        uuid.UUID
}
