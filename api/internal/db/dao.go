package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// This file contains the go structs that mimic the database tables

type Employer struct {
	ClientID           string      `db:"client_id"`
	OnboardStatus      string      `db:"onboard_status"`
	OnboardAdmin       pgtype.Text `db:"onboard_admin"`
	OnboardSecretToken pgtype.Text `db:"onboard_secret_token"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	OnboardEmailID pgtype.Int8 `db:"onboard_email_id"`
}

type EmailState string

const (
	EmailStatePending   EmailState = "PENDING"
	EmailStateProcessed EmailState = "PROCESSED"
)

type Email struct {
	ID            int64      `db:"id"`
	EmailFrom     string     `db:"email_from"`
	EmailTo       []string   `db:"email_to"`
	EmailCC       []string   `db:"email_cc"`
	EmailBCC      []string   `db:"email_bcc"`
	EmailSubject  string     `db:"email_subject"`
	EmailHTMLBody string     `db:"email_html_body"`
	EmailTextBody string     `db:"email_text_body"`
	EmailState    EmailState `db:"email_state"`

	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"`
}
