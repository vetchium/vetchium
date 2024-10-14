package db

import (
	"time"

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

// Structs
type HubUser struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type Email struct {
	ID            int64              `db:"id"`
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
	ID                 int64         `db:"id"`
	ClientIDType       ClientIDType  `db:"client_id_type"`
	EmployerState      EmployerState `db:"employer_state"`
	OnboardAdminEmail  string        `db:"onboard_admin_email"`
	OnboardSecretToken pgtype.Text   `db:"onboard_secret_token"`
	OnboardEmailID     pgtype.Int8   `db:"onboard_email_id"`
	CreatedAt          time.Time     `db:"created_at"`
}

type Domain struct {
	ID          int64       `db:"id"`
	DomainName  string      `db:"domain_name"`
	DomainState DomainState `db:"domain_state"`
	EmployerID  int64       `db:"employer_id"`
	CreatedAt   time.Time   `db:"created_at"`
}

type OrgUser struct {
	ID           int64       `db:"id"`
	Email        string      `db:"email"`
	PasswordHash string      `db:"password_hash"`
	OrgUserRole  OrgUserRole `db:"org_user_role"`
	EmployerID   int64       `db:"employer_id"`
	CreatedAt    time.Time   `db:"created_at"`
}
