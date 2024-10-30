package db

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

type Email struct {
	EmailKey      uuid.UUID          `db:"email_key"`
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
