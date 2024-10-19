package db

// This file contains the objects that should match the database schema.

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Enums
type EmailState string

const (
	EmailStatePending   EmailState = "PENDING"
	EmailStateProcessed EmailState = "PROCESSED"
)

type ClientIDType string

const (
	DomainClientIDType ClientIDType = "DOMAIN"
)

type EmployerState string

const (
	OnboardPendingEmployerState EmployerState = "ONBOARD_PENDING"
	OnboardedEmployerState      EmployerState = "ONBOARDED"
	DeboardedEmployerState      EmployerState = "DEBOARDED"
)

type DomainState string

const (
	VerifiedDomainState  DomainState = "VERIFIED"
	DeboardedDomainState DomainState = "DEBOARDED"
)

type OrgUserRole string

const (
	AdminOrgUserRole       OrgUserRole = "ADMIN"
	RecruiterOrgUserRole   OrgUserRole = "RECRUITER"
	InterviewerOrgUserRole OrgUserRole = "INTERVIEWER"
)

type OrgUserState string

const (
	ActiveOrgUserState OrgUserState = "ACTIVE"
	LockedOrgUserState OrgUserState = "LOCKED"
)

// Structs
type HubUser struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type Email struct {
	ID            uuid.UUID          `db:"id"`
	EmailFrom     string             `db:"email_from"`
	EmailTo       []string           `db:"email_to"`
	EmailCC       []string           `db:"email_cc"`
	EmailBCC      []string           `db:"email_bcc"`
	EmailSubject  string             `db:"email_subject"`
	EmailHTMLBody string             `db:"email_html_body"`
	EmailTextBody string             `db:"email_text_body"`
	EmailState    EmailState         `db:"email_state"`
	CreatedAt     time.Time          `db:"created_at"`
	ProcessedAt   pgtype.Timestamptz `db:"processed_at"`
}

type Employer struct {
	ID                 uuid.UUID          `db:"id"`
	ClientIDType       ClientIDType       `db:"client_id_type"`
	EmployerState      EmployerState      `db:"employer_state"`
	OnboardAdminEmail  string             `db:"onboard_admin_email"`
	OnboardSecretToken pgtype.Text        `db:"onboard_secret_token"`
	TokenValidTill     pgtype.Timestamptz `db:"token_valid_till"`
	OnboardEmailID     uuid.UUID          `db:"onboard_email_id"`
	CreatedAt          time.Time          `db:"created_at"`
}

type Domain struct {
	ID          uuid.UUID   `db:"id"`
	DomainName  string      `db:"domain_name"`
	DomainState DomainState `db:"domain_state"`
	EmployerID  uuid.UUID   `db:"employer_id"`
	CreatedAt   time.Time   `db:"created_at"`
}

type OrgUser struct {
	ID           uuid.UUID    `db:"id"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	OrgUserRole  OrgUserRole  `db:"org_user_role"`
	OrgUserState OrgUserState `db:"org_user_state"`
	EmployerID   uuid.UUID    `db:"employer_id"`
	CreatedAt    time.Time    `db:"created_at"`
}

type TokenType string

const (
	UserSessionToken TokenType = "USER_SESSION"
	TGToken          TokenType = "TGT"
	EmailToken       TokenType = "EMAIL"
)

type OrgUserToken struct {
	Token          string    `db:"token"`
	OrgUserID      uuid.UUID `db:"org_user_id"`
	TokenValidTill time.Time `db:"token_valid_till"`
	TokenType      TokenType `db:"token_type"`
	CreatedAt      time.Time `db:"created_at"`
}
