package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
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
	org_user_roles,
	org_user_state,
	employer_id
)
VALUES ($1, $2, $3::org_user_roles[], $4, $5)
`
	_, err = tx.Exec(
		ctx,
		orgUserInsertQuery,
		adminEmailAddr,
		onboardReq.Password,
		[]string{string(vetchi.Admin)},
		vetchi.ActiveOrgUserState,
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

		p.log.Error("failed to query org user auth", "error", err)
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
		p.log.Error("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

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
		p.log.Error("failed to insert TGT", "error", err)
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
) (db.OrgUserTO, error) {
	query := `
SELECT
    ou.id,
    ou.email,
    ou.employer_id,
    ou.org_user_roles,
    ou.org_user_state
FROM
    org_user_tokens out1,
    org_user_tokens out2,
    org_users ou
WHERE
    out1.token = $1
    AND out1.token_type = 'EMAIL'
    AND out2.token = $2
    AND out2.token_type = 'TGT'
    AND ou.id = out1.org_user_id
    AND ou.id = out2.org_user_id
`

	var orgUser db.OrgUserTO
	var roles []string
	err := p.pool.QueryRow(
		ctx,
		query, tfaCode, tgt).Scan(
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

		p.log.Error("failed to query org user by token", "error", err)
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
    AND out1.token_type = $2
    AND ou.id = out1.org_user_id
`

	var orgUser db.OrgUserTO
	var roles []string
	err := p.pool.QueryRow(
		ctx, query, sessionToken, db.EmployerSessionToken).Scan(
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

		p.log.Error("failed to query org user", "error", err)
		return db.OrgUserTO{}, err
	}

	orgUser.OrgUserRoles, err = p.convertToOrgUserRoles(roles)
	if err != nil {
		return db.OrgUserTO{}, err
	}

	return orgUser, nil
}

func (p *PG) convertToOrgUserRoles(
	dbRoles []string,
) ([]vetchi.OrgUserRole, error) {
	var roles []vetchi.OrgUserRole
	for _, str := range dbRoles {
		role := vetchi.OrgUserRole(str)
		switch role {
		case vetchi.Admin,
			vetchi.CostCentersCRUD,
			vetchi.CostCentersViewer,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
			vetchi.OpeningsCRUD,
			vetchi.OpeningsViewer:
			roles = append(roles, role)
		default:
			p.log.Error("invalid role in the database", "role", str)
			return nil, fmt.Errorf("invalid role: %s", str)
		}
	}
	return roles, nil
}

func (p *PG) GetEmployerByID(
	ctx context.Context,
	employerID uuid.UUID,
) (db.Employer, error) {
	query := "SELECT * FROM employers WHERE id = $1::UUID"
	rows, err := p.pool.Query(ctx, query, employerID)
	if err != nil {
		p.log.Error("failed to query employer", "error", err)
		return db.Employer{}, err
	}

	employer, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[db.Employer])
	if err != nil {
		p.log.Error("failed to collect one row", "error", err)
		return db.Employer{}, err
	}

	return employer, nil
}
