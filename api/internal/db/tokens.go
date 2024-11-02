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

	EmployerInviteToken TokenType = "EMPLOYER_INVITE"
)

// OrgUserTokenTO should be used to read from the database and NOT to
// write to the database. Use one of the TokenReq struct for Writes.
type OrgUserTokenTO struct {
	Token          string    `db:"token"`
	OrgUserID      uuid.UUID `db:"org_user_id"`
	TokenValidTill time.Time `db:"token_valid_till"`
	TokenType      TokenType `db:"token_type"`
	CreatedAt      time.Time `db:"created_at"`
}

type TokenReq struct {
	Token            string
	TokenType        TokenType
	ValidityDuration time.Duration
	OrgUserID        uuid.UUID
}
