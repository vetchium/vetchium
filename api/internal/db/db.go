package db

import (
	"context"
)

type DB interface {
	GetEmployer(ctx context.Context, clientID string) (Employer, error)
}
