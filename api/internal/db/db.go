package db

import (
	"context"
	"errors"
)

type DB interface {
	// Used by hermione and granger

	// Used by hermione
	InitEmployerAndDomain(
		ctx context.Context,
		employer Employer,
		domain Domain,
	) error
	GetEmployer(ctx context.Context, clientID string) (Employer, error)
	OnboardAdmin(
		ctx context.Context,
		domainName, password, token string,
	) error

	// Used by granger
	CleanOldOnboardTokens(ctx context.Context) error
	CreateOnboardEmail(
		ctx context.Context,
		employerID int64,
		onboardSecretToken string,
		tokenValidMins float64,
		email Email,
	) error
	GetOldestUnsentEmails(ctx context.Context) ([]Email, error)
	UpdateEmailState(ctx context.Context, emailID int64, state EmailState) error
	DeQOnboard(
		ctx context.Context,
	) (employerID int64, adminEmailAddr, domainName string, ok bool, err error)
}

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer           = errors.New("employer not found")
	ErrOrgUserAlreadyExists = errors.New("org user already exists")
)
