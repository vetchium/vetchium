package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/vetchium/vetchium/typespec/hub"
)

// Enums
type EmailState string

const (
	EmailStatePending   EmailState = "PENDING"
	EmailStateProcessed EmailState = "PROCESSED"
)

type Email struct {
	EmailKey      uuid.UUID  `db:"email_key"`
	EmailFrom     string     `db:"email_from"`
	EmailTo       []string   `db:"email_to"`
	EmailCC       []string   `db:"email_cc"`
	EmailBCC      []string   `db:"email_bcc"`
	EmailSubject  string     `db:"email_subject"`
	EmailHTMLBody string     `db:"email_html_body"`
	EmailTextBody string     `db:"email_text_body"`
	EmailState    EmailState `db:"email_state"`
	CreatedAt     time.Time  `db:"created_at"`
	ProcessedAt   time.Time  `db:"processed_at"`
}

type EmailStateChange struct {
	EmailDBKey uuid.UUID
	EmailState EmailState
}

type OnboardEmailInfo struct {
	EmployerID         uuid.UUID
	OnboardSecretToken string
	TokenValidMins     float64
	Email              Email
}

type ApplicationMailInfo struct {
	Employer EmployerMailInfo
	HubUser  HubUserMailInfo
	Opening  OpeningMailInfo
}

type OpeningMailInfo struct {
	OpeningID string
	Title     string
}

type HubUserMailInfo struct {
	HubUserID         uuid.UUID
	State             hub.HubUserState
	FullName          string
	Handle            string
	Email             string
	PreferredLanguage string
}

type EmployerMailInfo struct {
	EmployerID    uuid.UUID
	CompanyName   string
	PrimaryDomain string
}
