package postgres

import (
	"context"
	"errors"
	"log/slog"
	"strings"

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
SELECT
	e.id,
	e.client_id_type,
	e.employer_state,
	e.onboard_admin_email,
	e.onboard_secret_token,
	e.token_valid_till,
	e.created_at
FROM employers e, domains d
WHERE e.id = d.employer_id AND d.domain_name = $1
`

	var employer db.Employer
	err := p.pool.QueryRow(ctx, query, clientID).Scan(
		&employer.ID,
		&employer.ClientIDType,
		&employer.EmployerState,
		&employer.OnboardAdminEmail,
		&employer.OnboardSecretToken,
		&employer.TokenValidTill,
		&employer.CreatedAt,
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

func (p *PG) InitEmployerAndDomain(
	ctx context.Context,
	employer db.Employer,
	domain db.Domain,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	employerInsertQuery := `
INSERT INTO employers (
	client_id_type,
	employer_state,
	onboard_admin_email,
	onboard_secret_token
)
VALUES ($1, $2, $3, $4)
RETURNING id
`
	var employerID int64
	err = tx.QueryRow(
		ctx,
		employerInsertQuery,
		employer.ClientIDType,
		employer.EmployerState,
		employer.OnboardAdminEmail,
		nil,
	).Scan(&employerID)
	if err != nil {
		p.log.Error("failed to insert employer", "error", err)
		return err
	}

	domainInsertQuery := `
INSERT INTO domains (domain_name, domain_state, employer_id, created_at)
VALUES ($1, $2, $3, NOW())
`
	_, err = tx.Exec(
		ctx,
		domainInsertQuery,
		domain.DomainName,
		domain.DomainState,
		employerID,
	)
	if err != nil {
		p.log.Error("failed to insert domain", "error", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) WhomToOnboardInvite(
	ctx context.Context,
) (employerID int64, adminEmailAddr, domainName string, err error) {
	query := `
SELECT e.id, e.onboard_admin_email, d.domain_name
FROM employers e, domains d
WHERE e.employer_state = $1 AND e.onboard_secret_token IS NULL
ORDER BY e.created_at ASC
LIMIT 1
`
	err = p.pool.QueryRow(ctx, query, db.OnboardPendingEmployerState).
		Scan(&employerID, &adminEmailAddr, &domainName)
	if err != nil {
		p.log.Error("failed to query employers", "error", err)
		return 0, "", "", err
	}

	return employerID, adminEmailAddr, domainName, nil
}

func (p *PG) CreateOnboardEmail(
	ctx context.Context,
	employerID int64,
	onboardSecretToken string,
	tokenValidMins float64,
	email db.Email,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(),
		`
INSERT INTO emails (
	email_from,
	email_to,
	email_subject,
	email_html_body,
	email_text_body,
	email_state
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
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
UPDATE employers
SET
	onboard_email_id = $1,
	onboard_secret_token = $2, 
	token_valid_till = NOW() + interval '1 minute' * $3
WHERE id = $4
`,
		email.ID,
		onboardSecretToken,
		tokenValidMins,
		employerID,
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

func (p *PG) GetOldestUnsentEmails(ctx context.Context) ([]db.Email, error) {
	query := `
SELECT
	id,
	email_from,
	email_to,
	email_subject,
	email_html_body,
	email_text_body,
	email_state
FROM emails
WHERE email_state = $1
ORDER BY created_at ASC
LIMIT 10
`
	rows, err := p.pool.Query(ctx, query, db.EmailStatePending)
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

func (p *PG) UpdateEmailState(
	ctx context.Context,
	emailID int64,
	state db.EmailState,
) error {
	query := `
UPDATE emails
SET email_state = $1, processed_at = NOW()
WHERE id = $2
`
	_, err := p.pool.Exec(ctx, query, state, emailID)
	if err != nil {
		p.log.Error("failed to update email state", "error", err)
		return err
	}

	return nil
}

func (p *PG) OnboardAdmin(
	ctx context.Context,
	domainName, password, token string,
) error {
	employerQuery := `
SELECT e.id, e.onboard_admin_email
FROM employers e, domains d
WHERE e.onboard_secret_token = $1 AND d.domain_name = $2
`

	var employerID int64
	var adminEmailAddr string
	err := p.pool.QueryRow(ctx, employerQuery, token, domainName).
		Scan(&employerID, &adminEmailAddr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoEmployer
		}

		p.log.Error("failed to query employers", "error", err)
		return err
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	orgUserInsertQuery := `
INSERT INTO org_users (email, password_hash, org_user_role, employer_id)
VALUES ($1, $2, $3, $4)
`
	_, err = tx.Exec(
		ctx,
		orgUserInsertQuery,
		adminEmailAddr,
		password,
		db.AdminOrgUserRole,
		employerID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return db.ErrOrgUserAlreadyExists
		}

		p.log.Error("failed to insert org user", "error", err)
		return err
	}

	employerUpdateQuery := `
UPDATE employers
SET employer_state = $1, onboard_secret_token = NULL, token_valid_till = NULL
WHERE id = $2
`
	_, err = tx.Exec(
		ctx,
		employerUpdateQuery,
		db.OnboardedEmployerState,
		employerID,
	)
	if err != nil {
		p.log.Error("failed to update employer", "error", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) CleanOldOnboardTokens(ctx context.Context) error {
	// We hope that NTP will not be broken on the DB server and time
	// will not be set to future.
	query := `
UPDATE employers	
SET onboard_secret_token = NULL
WHERE onboard_secret_token IS NOT NULL AND token_valid_till < NOW()
`
	_, err := p.pool.Exec(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}

		p.log.Error("failed to clean old onboard tokens", "error", err)
		return err
	}

	return nil
}
