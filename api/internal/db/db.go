package db

import (
	"context"
	"errors"
)

type DB interface {
	GetEmployer(ctx context.Context, clientID string) (Employer, error)
	CreateEmployer(ctx context.Context, employer Employer) error
	GetUnmailedOnboardPendingEmployers() ([]Employer, error)
	CreateOnboardEmail(employer Employer, email Email) error
}

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer = errors.New("employer not found")
)
