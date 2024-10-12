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
SELECT 	client_id, onboard_status, onboard_admin, 
		onboard_secret_token, onboard_email_id, 
		created_at, updated_at
FROM employers 
WHERE client_id = $1`

	var employer db.Employer

	err := p.pool.QueryRow(ctx, query, clientID).
		Scan(
			&employer.ClientID,
			&employer.OnboardStatus,
			&employer.OnboardAdmin,
			&employer.OnboardSecretToken,
			&employer.OnboardEmailID,
			&employer.CreatedAt,
			&employer.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Employer{}, db.ErrNoEmployer
		}

		p.log.Error("failed to get employer", "error", err)
		return db.Employer{}, err
	}

	return employer, nil
}

func (p *PG) CreateEmployer(
	ctx context.Context,
	employer db.Employer,
) error {
	query := `
INSERT INTO employers	(client_id, onboard_status, onboard_admin, 
						onboard_secret_token, onboard_email_id)
VALUES ($1, $2, $3, $4, $5)
`
	_, err := p.pool.Exec(
		ctx,
		query,
		employer.ClientID,
		employer.OnboardStatus,
		employer.OnboardAdmin,
		nil,
		nil,
	)
	return err
}

func (p *PG) GetUnmailedOnboardPendingEmployers() ([]db.Employer, error) {
	query := `
SELECT 	client_id, onboard_status, onboard_admin, 
		created_at, updated_at, onboard_email_id
FROM 	employers 
WHERE 	onboard_status = 'DOMAIN_VERIFIED_ONBOARD_PENDING'
		AND onboard_email_id IS NULL
`
	rows, err := p.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employers []db.Employer
	for rows.Next() {
		var employer db.Employer
		err := rows.Scan(
			&employer.ClientID,
			&employer.OnboardStatus,
			&employer.OnboardAdmin,
			&employer.CreatedAt,
			&employer.UpdatedAt,
			&employer.OnboardEmailID,
		)
		if err != nil {
			return nil, err
		}

		employers = append(employers, employer)
	}

	return employers, nil
}

func (p *PG) CreateOnboardEmail(employer db.Employer, email db.Email) error {
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(),
		`
INSERT INTO emails 	(email_from, email_to, email_subject, 
					email_html_body, email_text_body, email_state)
VALUES 				($1, $2, $3, $4, $5, $6)
RETURNING 			id
`,
		email.EmailFrom,
		email.EmailTo,
		email.EmailSubject,
		email.EmailHTMLBody,
		email.EmailTextBody,
		db.EmailStatePending,
	).Scan(&email.ID)
	if err != nil {
		p.log.Error("failed to create onboard email", "error", err)
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
UPDATE 	employers 
SET 	onboard_email_id = $1, onboard_secret_token = $2
WHERE 	client_id = $3
`,
		email.ID,
		employer.OnboardSecretToken,
		employer.ClientID,
	)
	if err != nil {
		p.log.Error("failed to update employer", "error", err)
		return err
	}
	// TODO: Ensure that the update query above correctly updates updated_at

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetOldestUnsentEmails() ([]db.Email, error) {
	query := `
SELECT 	id, email_from, email_to, email_subject, 
		email_html_body, email_text_body, email_state
FROM 	emails
WHERE 	email_state = $1
ORDER BY created_at ASC
LIMIT 10
`
	rows, err := p.pool.Query(context.Background(), query, db.EmailStatePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []db.Email
	for rows.Next() {
		var email db.Email
		err := rows.Scan(
			&email.ID,
			&email.EmailFrom,
			&email.EmailTo,
			&email.EmailSubject,
			&email.EmailHTMLBody,
			&email.EmailTextBody,
			&email.EmailState,
		)
		if err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (p *PG) UpdateEmailState(emailID int64, state db.EmailState) error {
	query := `UPDATE emails SET email_state = $1 WHERE id = $2`
	_, err := p.pool.Exec(context.Background(), query, state, emailID)
	return err
}
