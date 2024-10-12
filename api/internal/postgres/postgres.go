package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/psankar/vetchi/api/internal/db"
)

type PG struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func New(connStr string, logger *slog.Logger) (*PG, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	cdb := PG{pool: pool, log: logger}
	return &cdb, nil
}

func (p *PG) GetEmployer(
	ctx context.Context,
	clientID string,
) (db.Employer, error) {
	query := `
SELECT client_id, onboard_status, onboarding_admin, onboarding_email_sent_at, created_at, updated_at
FROM employers 
WHERE client_id = $1`

	var employer db.Employer

	err := p.pool.QueryRow(ctx, query, clientID).
		Scan(
			&employer.ClientID,
			&employer.OnboardStatus,
			&employer.OnboardingAdmin,
			&employer.OnboardingEmailSentAt,
			&employer.CreatedAt,
			&employer.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Employer{}, db.ErrNoEmployer
		}

		return db.Employer{}, err
	}

	return employer, nil
}

func (p *PG) CreateEmployer(
	ctx context.Context,
	employer db.Employer,
) error {
	query := `
INSERT INTO employers (client_id, onboard_status, onboarding_admin)
VALUES ($1, $2, $3)
`
	_, err := p.pool.Exec(
		ctx,
		query,
		employer.ClientID,
		employer.OnboardStatus,
		employer.OnboardingAdmin,
	)
	return err
}
