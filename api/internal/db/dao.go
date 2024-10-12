package db

import "time"

// This file contains the go structs that mimic the database tables

type Employer struct {
	ClientID              string     `db:"client_id"`
	OnboardStatus         string     `db:"onboard_status"`
	OnboardingAdmin       string     `db:"onboarding_admin"`
	OnboardingEmailSentAt *time.Time `db:"onboarding_email_sent_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
