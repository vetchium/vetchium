package postgres

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
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
	var employerID uuid.UUID
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

func (p *PG) DeQOnboard(
	ctx context.Context,
) (*db.OnboardInfo, error) {
	query := `
SELECT e.id, e.onboard_admin_email, d.domain_name
FROM employers e, domains d
WHERE 	e.employer_state = $1 AND
		e.onboard_secret_token IS NULL AND
		d.employer_id = e.id
ORDER BY e.created_at ASC
LIMIT 1
`

	var onboardInfo db.OnboardInfo
	err := p.pool.QueryRow(ctx, query, db.OnboardPendingEmployerState).
		Scan(
			&onboardInfo.EmployerID,
			&onboardInfo.AdminEmailAddr,
			&onboardInfo.DomainName,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		p.log.Error("failed to query employers", "error", err)
		return nil, err
	}

	return &onboardInfo, nil
}

func (p *PG) CreateOnboardEmail(
	ctx context.Context,
	onboardEmailInfo db.OnboardEmailInfo,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	var emailTableKey uuid.UUID
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
RETURNING email_key
`,
		onboardEmailInfo.Email.EmailFrom,
		onboardEmailInfo.Email.EmailTo,
		onboardEmailInfo.Email.EmailSubject,
		onboardEmailInfo.Email.EmailHTMLBody,
		onboardEmailInfo.Email.EmailTextBody,
		db.EmailStatePending,
	).Scan(&emailTableKey)
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
		emailTableKey,
		onboardEmailInfo.OnboardSecretToken,
		onboardEmailInfo.TokenValidMins,
		onboardEmailInfo.EmployerID,
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
	email_key,
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
			&email.EmailKey,
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
	emailStateChange db.EmailStateChange,
) error {
	query := `
UPDATE emails
SET email_state = $1, processed_at = NOW()
WHERE email_key = $2
`
	_, err := p.pool.Exec(
		ctx,
		query,
		emailStateChange.EmailState,
		emailStateChange.EmailDBKey,
	)
	if err != nil {
		p.log.Error("failed to update email state", "error", err)
		return err
	}

	return nil
}

func (p *PG) OnboardAdmin(
	ctx context.Context,
	onboardReq db.OnboardReq,
) error {
	employerQuery := `
SELECT e.id, e.onboard_admin_email
FROM employers e, domains d
WHERE
	e.onboard_secret_token = $1 AND
	e.token_valid_till > NOW() AND
	d.domain_name = $2 AND
	d.employer_id = e.id AND
	e.employer_state = $3
`

	var employerID uuid.UUID
	var adminEmailAddr string
	err := p.pool.QueryRow(
		ctx,
		employerQuery,
		onboardReq.Token,
		onboardReq.DomainName,
		db.OnboardPendingEmployerState,
	).Scan(&employerID, &adminEmailAddr)
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
INSERT INTO org_users (
	email,
	password_hash,
	org_user_role,
	org_user_state,
	employer_id
)
VALUES ($1, $2, $3, $4, $5)
`
	_, err = tx.Exec(
		ctx,
		orgUserInsertQuery,
		adminEmailAddr,
		onboardReq.Password,
		db.AdminOrgUserRole,
		db.ActiveOrgUserState,
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
UPDATE employers SET employer_state = $1 WHERE id = $2
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

func (p *PG) PruneTokens(ctx context.Context) error {
	// We hope that NTP will not be broken on the DB server and time
	// will not be set to future in the db server.

	queries := []string{
		`DELETE FROM org_user_tokens WHERE token_valid_till < NOW()`,
	}

	for _, q := range queries {
		_, err := p.pool.Exec(ctx, q)
		if err != nil {
			if err == pgx.ErrNoRows {
				// No obsolete tokens to delete
				continue
			}

			p.log.Error("failed to execute query", "error", err)
			return err
		}
	}

	return nil
}

func (p *PG) GetOrgUserAuth(
	ctx context.Context,
	orgUserCreds db.OrgUserCreds,
) (db.OrgUserAuth, error) {
	query := `
SELECT
	ou.id,
	ou.email,
	ou.employer_id,
	ou.org_user_role,
	ou.password_hash,
	e.employer_state,
	ou.org_user_state
FROM org_users ou, employers e, domains d
WHERE 	ou.email = $1 AND
		ou.employer_id = e.id AND
		e.id = d.employer_id AND
		d.domain_name = $2
`

	var orgUserAuth db.OrgUserAuth
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUserCreds.Email,
		orgUserCreds.ClientID,
	).Scan(
		&orgUserAuth.OrgUserID,
		&orgUserAuth.OrgUserEmail,
		&orgUserAuth.EmployerID,
		&orgUserAuth.OrgUserRole,
		&orgUserAuth.PasswordHash,
		&orgUserAuth.EmployerState,
		&orgUserAuth.OrgUserState,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUserAuth{}, db.ErrNoOrgUser
		}

		p.log.Error("failed to query org user auth", "error", err)
		return db.OrgUserAuth{}, err
	}

	return orgUserAuth, nil
}

func (p *PG) InitEmployerTFA(
	ctx context.Context,
	employerTFA db.EmployerTFA,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, $3, $4)
`,
		employerTFA.TGToken.Token,
		employerTFA.TGToken.OrgUserID,
		employerTFA.TGToken.TokenValidTill,
		employerTFA.TGToken.TokenType,
	)
	if err != nil {
		p.log.Error("failed to insert TGT", "error", err)
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, $3, $4)
`,
		employerTFA.EmailToken.Token,
		employerTFA.EmailToken.OrgUserID,
		employerTFA.EmailToken.TokenValidTill,
		employerTFA.EmailToken.TokenType,
	)
	if err != nil {
		p.log.Error("failed to insert EMAILToken", "error", err)
		return err
	}

	_, err = tx.Exec(
		ctx,
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
`,
		employerTFA.Email.EmailFrom,
		employerTFA.Email.EmailTo,
		employerTFA.Email.EmailSubject,
		employerTFA.Email.EmailHTMLBody,
		employerTFA.Email.EmailTextBody,
		employerTFA.Email.EmailState,
	)
	if err != nil {
		p.log.Error("failed to insert Email", "error", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetOrgUserByToken(
	ctx context.Context,
	tfaCode, tgt string,
) (db.OrgUser, error) {
	query := `
SELECT
	ou.id,
	ou.email,
	ou.employer_id,
	ou.org_user_role,
	ou.org_user_state
FROM org_user_tokens out1, org_user_tokens out2, org_users ou
WHERE
	out1.token = $1 AND out1.token_type = 'EMAIL' AND
	out2.token = $2 AND out2.token_type = 'TGT' AND
	ou.id = out1.org_user_id AND
	ou.id = out2.org_user_id
`

	var orgUser db.OrgUser
	err := p.pool.QueryRow(
		ctx,
		query, tfaCode, tgt).Scan(
		&orgUser.ID,
		&orgUser.Email,
		&orgUser.EmployerID,
		&orgUser.OrgUserRole,
		&orgUser.OrgUserState,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUser{}, db.ErrNoOrgUser
		}

		p.log.Error("failed to query org user by token", "error", err)
		return db.OrgUser{}, err
	}

	return orgUser, nil
}

func (p *PG) CreateOrgUserToken(
	ctx context.Context,
	orgUserToken db.OrgUserToken,
) error {
	query := `
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, $3, $4)
`
	_, err := p.pool.Exec(
		ctx,
		query,
		orgUserToken.Token,
		orgUserToken.OrgUserID,
		orgUserToken.TokenValidTill,
		orgUserToken.TokenType,
	)
	if err != nil {
		// TODO: Check if the error is due to duplicate key value
		// and if so retry with a different token
		p.log.Error("failed to create org user token", "error", err)
		return err
	}

	return nil
}

func (p *PG) CreateCostCenter(
	ctx context.Context,
	costCenterReq db.CostCenterReq,
) (uuid.UUID, error) {
	query := `
INSERT INTO org_cost_centers (cost_center_name, notes, employer_id)
VALUES ($1, $2, $3)
RETURNING id
`
	var costCenterID uuid.UUID
	err := p.pool.QueryRow(
		ctx, query,
		costCenterReq.Name,
		costCenterReq.Notes,
		costCenterReq.OrgUserID,
	).Scan(&costCenterID)
	if err != nil {
		p.log.Error("failed to create cost center", "error", err)
		return uuid.UUID{}, err
	}

	return costCenterID, nil
}

func (p *PG) AuthOrgUser(
	ctx context.Context,
	sessionToken string,
) (db.OrgUser, error) {
	query := `
SELECT
	ou.id,
	ou.email,
	ou.employer_id,
	ou.org_user_role,
	ou.org_user_state
FROM org_user_tokens out1, org_users ou
WHERE out1.token = $1 AND out1.token_type = 'SESSION' AND
	ou.id = out1.org_user_id
`

	var orgUser db.OrgUser
	err := p.pool.QueryRow(
		ctx, query, sessionToken).Scan(
		&orgUser.ID,
		&orgUser.Email,
		&orgUser.EmployerID,
		&orgUser.OrgUserRole,
		&orgUser.OrgUserState,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUser{}, db.ErrNoOrgUser
		}

		p.log.Error("failed to query org user", "error", err)
		return db.OrgUser{}, err
	}

	return orgUser, nil
}
