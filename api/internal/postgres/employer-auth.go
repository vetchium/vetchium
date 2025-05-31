package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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
FROM
    employers e,
    domains d
WHERE
    e.id = d.employer_id
    AND d.domain_name = $1
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

		p.log.Err("failed to get employer", "error", err)
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
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	employerInsertQuery := `
INSERT INTO employers (
	client_id_type,
	employer_state,
	company_name,
	onboard_admin_email,
	onboard_secret_token
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`
	var employerID uuid.UUID
	err = tx.QueryRow(
		ctx,
		employerInsertQuery,
		employer.ClientIDType,
		employer.EmployerState,
		domain.DomainName,
		employer.OnboardAdminEmail,
		nil,
	).Scan(&employerID)
	if err != nil {
		p.log.Err("failed to insert employer", "error", err)
		return err
	}

	// Single query to either update an existing unverified domain or insert a new one
	domainUpsertQuery := `
INSERT INTO domains (domain_name, domain_state, employer_id, created_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (domain_name) DO UPDATE
SET domain_state = EXCLUDED.domain_state,
    employer_id = EXCLUDED.employer_id
WHERE domains.domain_state = $4
AND domains.employer_id IS NULL
RETURNING id`

	var domainID uuid.UUID
	err = tx.QueryRow(
		ctx,
		domainUpsertQuery,
		domain.DomainName,
		db.VerifiedDomainState,
		employerID,
		db.UnverifiedDomainState,
	).Scan(&domainID)
	if err != nil {
		p.log.Err("failed to upsert domain", "error", err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) DeQOnboard(
	ctx context.Context,
) (*db.OnboardInfo, error) {
	query := `
SELECT
    e.id,
    e.onboard_admin_email,
    d.domain_name
FROM
    employers e,
    domains d
WHERE
    e.employer_state = $1
    AND e.onboard_secret_token IS NULL
    AND d.employer_id = e.id
ORDER BY
    e.created_at ASC
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

		p.log.Err("failed to query employers", "error", err)
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
		p.log.Err("failed to begin transaction", "error", err)
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
		p.log.Err("failed to create onboard email", "error", err)
		return err
	}

	_, err = tx.Exec(
		context.Background(),
		`
UPDATE
    employers
SET
    onboard_email_id = $1,
    onboard_secret_token = $2,
    token_valid_till = (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute'))
WHERE
    id = $4
`,
		emailTableKey,
		onboardEmailInfo.OnboardSecretToken,
		onboardEmailInfo.TokenValidMins,
		onboardEmailInfo.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to update employer", "error", err)
		return err
	}
	// TODO: Ensure that the update query above correctly updates updated_at

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) OnboardAdmin(
	ctx context.Context,
	onboardReq db.OnboardReq,
) error {
	employerQuery := `
SELECT
    e.id,
    e.onboard_admin_email
FROM
    employers e,
    domains d
WHERE
    e.onboard_secret_token = $1
    AND e.token_valid_till > NOW()
    AND d.domain_name = $2
    AND d.employer_id = e.id
    AND e.employer_state = $3
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

		p.log.Err("failed to query employers", "error", err)
		return err
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	orgUserInsertQuery := `
INSERT INTO org_users (
	name,
	email,
	password_hash,
	org_user_roles,
	org_user_state,
	employer_id
)
VALUES ($1, $2, $3, $4::org_user_roles[], $5, $6)
`
	_, err = tx.Exec(
		ctx,
		orgUserInsertQuery,

		// During onboarding, we will use the admin email as the name.
		// The admin can change this later.
		adminEmailAddr,

		adminEmailAddr,
		onboardReq.Password,
		[]string{string(common.Admin)},
		employer.ActiveOrgUserState,
		employerID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return db.ErrOrgUserAlreadyExists
		}

		p.log.Err("failed to insert org user", "error", err)
		return err
	}

	primaryDomainInsertQuery := `
INSERT INTO employer_primary_domains (employer_id, domain_id)
VALUES ($1, (SELECT id FROM domains WHERE domain_name = $2))
`
	_, err = tx.Exec(
		ctx,
		primaryDomainInsertQuery,
		employerID,
		onboardReq.DomainName,
	)
	if err != nil {
		p.log.Err("failed to insert primary domain", "error", err)
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
		p.log.Err("failed to update employer", "error", err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) PruneTokens(ctx context.Context) error {
	// We hope that NTP will not be broken on the DB server and time
	// will not be set to future in the db server.

	queries := []string{
		`DELETE FROM org_user_tokens WHERE token_valid_till < NOW()`,
		`DELETE FROM hub_user_tokens WHERE token_valid_till < NOW()`,
	}

	for _, q := range queries {
		_, err := p.pool.Exec(ctx, q)
		if err != nil {
			if err == pgx.ErrNoRows {
				// No obsolete tokens to delete
				continue
			}

			p.log.Err("failed to execute query", "error", err)
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
    ou.org_user_roles,
    ou.password_hash,
    e.employer_state,
    ou.org_user_state
FROM
    org_users ou,
    employers e,
    domains d
WHERE
    ou.email = $1
    AND ou.employer_id = e.id
    AND e.id = d.employer_id
    AND d.domain_name = $2
`

	var orgUserAuth db.OrgUserAuth
	var roles []string
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUserCreds.Email,
		orgUserCreds.ClientID,
	).Scan(
		&orgUserAuth.OrgUserID,
		&orgUserAuth.OrgUserEmail,
		&orgUserAuth.EmployerID,
		&roles,
		&orgUserAuth.PasswordHash,
		&orgUserAuth.EmployerState,
		&orgUserAuth.OrgUserState,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUserAuth{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to query org user auth", "error", err)
		return db.OrgUserAuth{}, err
	}

	orgUserAuth.OrgUserRoles, err = p.convertToOrgUserRoles(roles)
	if err != nil {
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
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		ctx,
		`
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`,
		employerTFA.TFAToken.Token,
		employerTFA.TFAToken.OrgUserID,
		employerTFA.TFAToken.ValidityDuration.Minutes(),
		db.EmployerTFAToken,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "org_user_tokens_pkey" {
			p.log.Err("duplicate token generated", "error", err)
			return err
		}
		p.log.Err("failed to insert TFA Token", "error", err)
		return err
	}

	_, err = tx.Exec(
		ctx,
		`
INSERT INTO org_user_tfa_codes(code, tfa_token)
VALUES ($1, $2)
`,
		employerTFA.TFACode,
		employerTFA.TFAToken.Token,
	)
	if err != nil {
		p.log.Err("failed to insert TFA code", "error", err)
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
		p.log.Err("failed to insert Email", "error", err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetOrgUserByTFACreds(
	ctx context.Context,
	tfaCode, tfaToken string,
) (db.OrgUserTO, error) {
	query := `
SELECT
    ou.id,
    ou.email,
    ou.employer_id,
    ou.org_user_roles,
    ou.org_user_state
FROM
	org_user_tfa_codes oc,
	org_user_tokens ot,
	org_users ou
WHERE
	oc.code = $1
	AND ot.token = $2
	AND oc.tfa_token = ot.token
	AND ot.token_type = $3
	AND ou.id = ot.org_user_id
`

	var orgUser db.OrgUserTO
	var roles []string
	err := p.pool.QueryRow(
		ctx,
		query,
		tfaCode,
		tfaToken,
		db.EmployerTFAToken,
	).Scan(
		&orgUser.ID,
		&orgUser.Email,
		&orgUser.EmployerID,
		&roles,
		&orgUser.OrgUserState,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUserTO{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to query org user by token", "error", err)
		return db.OrgUserTO{}, err
	}

	orgUser.OrgUserRoles, err = p.convertToOrgUserRoles(roles)
	if err != nil {
		return db.OrgUserTO{}, err
	}

	return orgUser, nil
}

func (p *PG) AuthOrgUser(
	ctx context.Context,
	sessionToken string,
) (db.OrgUserTO, error) {
	query := `
SELECT
    ou.id,
    ou.email,
    ou.employer_id,
    ou.org_user_roles,
    ou.org_user_state,
    ou.created_at
FROM
    org_user_tokens out1,
    org_users ou
WHERE
    out1.token = $1
    AND (out1.token_type = $2 OR out1.token_type = $3)
    AND ou.id = out1.org_user_id
`

	var orgUser db.OrgUserTO
	var roles []string
	err := p.pool.QueryRow(
		ctx,
		query,
		sessionToken,
		db.EmployerSessionToken,
		db.EmployerLTSToken,
	).Scan(
		&orgUser.ID,
		&orgUser.Email,
		&orgUser.EmployerID,
		&roles,
		&orgUser.OrgUserState,
		&orgUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.OrgUserTO{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to query org user", "error", err)
		return db.OrgUserTO{}, err
	}

	orgUser.OrgUserRoles, err = p.convertToOrgUserRoles(roles)
	if err != nil {
		return db.OrgUserTO{}, err
	}

	return orgUser, nil
}

func (p *PG) GetEmployerByID(
	ctx context.Context,
	employerID uuid.UUID,
) (db.Employer, error) {
	query := "SELECT * FROM employers WHERE id = $1::UUID"
	rows, err := p.pool.Query(ctx, query, employerID)
	if err != nil {
		p.log.Err("failed to query employer", "error", err)
		return db.Employer{}, err
	}

	employer, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[db.Employer])
	if err != nil {
		p.log.Err("failed to collect one row", "error", err)
		return db.Employer{}, err
	}

	return employer, nil
}

func (p *PG) InitEmployerPasswordReset(
	ctx context.Context,
	initPasswordResetReq db.EmployerInitPasswordReset,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	tokensQuery := `
INSERT INTO org_user_tokens(token, org_user_id, token_valid_till, token_type)
VALUES ($1, $2, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')), $4)
`
	_, err = tx.Exec(
		ctx,
		tokensQuery,
		initPasswordResetReq.Token,
		initPasswordResetReq.OrgUserID,
		initPasswordResetReq.ValidityDuration.Minutes(),
		db.EmployerResetPasswordToken,
	)
	if err != nil {
		p.log.Err("failed to insert password reset token", "error", err)
		return err
	}

	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		initPasswordResetReq.Email.EmailFrom,
		initPasswordResetReq.Email.EmailTo,
		initPasswordResetReq.Email.EmailSubject,
		initPasswordResetReq.Email.EmailHTMLBody,
		initPasswordResetReq.Email.EmailTextBody,
		initPasswordResetReq.Email.EmailState,
	)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) ResetEmployerPassword(
	ctx context.Context,
	employerPasswordReset db.EmployerPasswordReset,
) error {
	query := `
WITH token_info AS (
    SELECT org_user_id, token
    FROM org_user_tokens
    WHERE token = $2
),
password_update AS (
    UPDATE org_users
    SET password_hash = $1
	WHERE id = (SELECT org_user_id FROM token_info)
    RETURNING id
)
DELETE FROM org_user_tokens
WHERE token = (
    SELECT token
    FROM token_info
)
AND EXISTS (SELECT 1 FROM password_update)
RETURNING (SELECT id FROM password_update)
`
	var orgUserID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		employerPasswordReset.PasswordHash,
		employerPasswordReset.Token,
	).Scan(&orgUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("invalid password reset token", "error", err)
			return db.ErrInvalidPasswordResetToken
		}

		p.log.Err("failed to reset password", "error", err)
		return err
	}

	p.log.Dbg("password reset", "orgUserID", orgUserID)

	return nil
}
