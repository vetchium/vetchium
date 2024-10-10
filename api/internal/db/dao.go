package db

// This file contains the go structs that mimic the database tables

type Employer struct {
	ClientID      string `db:"client_id"`
	OnboardStatus string `db:"onboard_status"`
}
