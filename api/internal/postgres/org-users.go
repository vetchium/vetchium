package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) AddOrgUser(
	ctx context.Context,
	addOrgUserReq db.AddOrgUserReq,
) (orgUserID uuid.UUID, err error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return uuid.UUID{}, err
	}

	defer tx.Rollback(ctx)

	orgUserQuery := `
INSERT INTO org_users (name, email, employer_id, org_user_roles, org_user_state)
	VALUES ($1, $2, $3::UUID, $4::org_user_roles[], $5)
RETURNING id
`
	row := tx.QueryRow(
		ctx,
		orgUserQuery,
		addOrgUserReq.Name,
		addOrgUserReq.Email,
		addOrgUserReq.EmployerID,
		addOrgUserReq.OrgUserRoles.StringArray(),
		addOrgUserReq.OrgUserState,
	)
	err = row.Scan(&orgUserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_email_employer_id" {
			return uuid.UUID{}, db.ErrOrgUserAlreadyExists
		}

		p.log.Err("failed to add org user", "error", err)
		return uuid.UUID{}, err
	}
	p.log.Dbg("org user added", "org_user_id", orgUserID)

	var tokenQuery = `
INSERT INTO org_user_invites(token, org_user_id, token_valid_till)
	VALUES ($1, $2::UUID, (NOW() AT TIME ZONE 'utc' + ($3 * INTERVAL '1 minute')))
RETURNING token
`
	var tokenKey string
	row = tx.QueryRow(
		ctx,
		tokenQuery,
		addOrgUserReq.InviteToken.Token,
		orgUserID,
		addOrgUserReq.InviteToken.ValidityDuration.Minutes(),
	)
	err = row.Scan(&tokenKey)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "org_user_invites_pkey" {
			// Unlikely to happen but still handle it
			p.log.Err("duplicate token generated", "error", err)
			return uuid.UUID{}, err
		}
		p.log.Err("failed to add org user token", "error", err)
		return uuid.UUID{}, err
	}
	p.log.Dbg("org user token added", "token_key", tokenKey)

	var emailQuery = `
INSERT INTO emails(email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
	VALUES ($1, $2, $3, $4, $5, $6)
RETURNING email_key
`
	var emailKey uuid.UUID
	row = tx.QueryRow(
		ctx,
		emailQuery,
		addOrgUserReq.InviteMail.EmailFrom,
		addOrgUserReq.InviteMail.EmailTo,
		addOrgUserReq.InviteMail.EmailSubject,
		addOrgUserReq.InviteMail.EmailHTMLBody,
		addOrgUserReq.InviteMail.EmailTextBody,
		addOrgUserReq.InviteMail.EmailState,
	)
	err = row.Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to add invite email", "error", err)
		return uuid.UUID{}, err
	}
	p.log.Dbg("invite email added", "email_key", emailKey)

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return uuid.UUID{}, err
	}

	return orgUserID, nil
}

func (p *PG) DisableOrgUser(
	ctx context.Context,
	disableOrgUserRequest employer.DisableOrgUserRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	const (
		userNotFound    = "USER_NOT_FOUND"
		lastActiveAdmin = "LAST_ACTIVE_ADMIN"
		alreadyDisabled = "ALREADY_DISABLED"
	)
	query := `
WITH target_user AS (
    SELECT id, org_user_roles, org_user_state
    FROM org_users
    WHERE employer_id = $1
      AND email = $2
),
other_active_admins AS (
    SELECT 1
    FROM org_users
    WHERE employer_id = $1
      AND id != (SELECT id FROM target_user)
      AND 'ADMIN' = ANY(org_user_roles::text[])
      AND org_user_state = $3 -- ACTIVE_ORG_USER
    LIMIT 1
),
updated_user AS (
    UPDATE org_users
    SET org_user_state = $4 -- DISABLED_ORG_USER
    WHERE id = (SELECT id FROM target_user)
      AND org_user_state != $4 -- DISABLED_ORG_USER
      AND (
          NOT 'ADMIN' = ANY((SELECT org_user_roles FROM target_user)::text[])
          OR EXISTS (SELECT 1 FROM other_active_admins)
      )
    RETURNING id
),
deleted_tokens AS (
    DELETE FROM org_user_tokens
    WHERE org_user_id = (SELECT id FROM updated_user)
)
SELECT 
    CASE 
        WHEN (SELECT id FROM target_user) IS NULL THEN $5
        WHEN (SELECT org_user_state FROM target_user) = $4 THEN $7
        WHEN 'ADMIN' = ANY((SELECT org_user_roles FROM target_user)::text[])
          AND NOT EXISTS (SELECT 1 FROM other_active_admins) THEN $6
        ELSE COALESCE((SELECT id FROM updated_user)::TEXT, $6)
    END AS result
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUser.EmployerID,
		disableOrgUserRequest.Email,
		employer.ActiveOrgUserState,
		employer.DisabledOrgUserState,
		userNotFound,
		lastActiveAdmin,
		alreadyDisabled,
	).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNoOrgUser
		}

		p.log.Err("failed to disable org user", "error", err)
		return err
	}

	switch result {
	case userNotFound:
		return db.ErrNoOrgUser
	case lastActiveAdmin:
		return db.ErrLastActiveAdmin
	case alreadyDisabled:
		return db.ErrOrgUserAlreadyDisabled
	default:
		orgUserID, err := uuid.Parse(result)
		if err != nil {
			p.log.Err("failed to parse org user id", "error", err)
			return err
		}

		p.log.Dbg("org user disabled", "org_user_id", orgUserID)
		return nil
	}
}

func (p *PG) EnableOrgUser(
	ctx context.Context,
	enableOrgUserReq db.EnableOrgUserReq,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}

	defer tx.Rollback(ctx)

	getOrgUserQuery := `
SELECT
    id,
    org_user_state
FROM
    org_users
WHERE
    email = $1
	AND employer_id = $2
`

	var orgUserID uuid.UUID
	var orgUserState employer.OrgUserState
	err = tx.QueryRow(ctx, getOrgUserQuery, enableOrgUserReq.Email, enableOrgUserReq.EmployerID).
		Scan(&orgUserID, &orgUserState)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNoOrgUser
		}

		p.log.Err("failed to get org user", "error", err)
		return err
	}

	if orgUserState != employer.DisabledOrgUserState {
		p.log.Dbg("org user not disabled", "org_user_id", orgUserID)
		return db.ErrOrgUserNotDisabled
	}

	updateOrgUserQuery := `
UPDATE
    org_users
SET
    org_user_state = $1 -- ADDED_ORG_USER
WHERE
    id = $2
    AND org_user_state = $3 -- DISABLED_ORG_USER
RETURNING id
`

	var updatedOrgUserID uuid.UUID
	err = tx.QueryRow(
		ctx,
		updateOrgUserQuery,
		employer.AddedOrgUserState,
		orgUserID,
		employer.DisabledOrgUserState,
	).Scan(&updatedOrgUserID)
	if err != nil {
		p.log.Err("failed to update org user", "error", err)
		return err
	}
	p.log.Dbg("org user set to Added state", "org_user_id", updatedOrgUserID)

	var emailQuery = `
INSERT INTO emails(email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
	VALUES ($1, $2, $3, $4, $5, $6)
RETURNING email_key
`
	var emailKey uuid.UUID
	err = tx.QueryRow(
		ctx,
		emailQuery,
		enableOrgUserReq.InviteMail.EmailFrom,
		enableOrgUserReq.InviteMail.EmailTo,
		enableOrgUserReq.InviteMail.EmailSubject,
		enableOrgUserReq.InviteMail.EmailHTMLBody,
		enableOrgUserReq.InviteMail.EmailTextBody,
		db.EmailStatePending,
	).Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to add invite email", "error", err)
		return err
	}
	p.log.Dbg("invite email added", "email_key", emailKey)

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) FilterOrgUsers(
	ctx context.Context,
	filterOrgUsersReq employer.FilterOrgUsersRequest,
) ([]employer.OrgUser, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	query := `
SELECT
    name,
    email,
    org_user_roles::text[],
    org_user_state
FROM
    org_users
WHERE
    employer_id = $1
    AND email > $2
    AND org_user_state = ANY ($3::org_user_states[])
    AND (email ILIKE $4
        OR name ILIKE $4)
ORDER BY
    email
LIMIT $5
`

	prefix := filterOrgUsersReq.Prefix + "%"
	rows, err := p.pool.Query(
		ctx,
		query,
		orgUser.EmployerID,
		filterOrgUsersReq.PaginationKey,
		filterOrgUsersReq.StatesAsStrings(),
		prefix,
		filterOrgUsersReq.Limit,
	)
	if err != nil {
		p.log.Err("failed to filter org users", "error", err)
		return nil, err
	}

	orgUsers, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[employer.OrgUser],
	)
	if err != nil {
		p.log.Err("failed to collect org users", "error", err)
		return nil, err
	}

	return orgUsers, nil
}

func (p *PG) SignupOrgUser(
	ctx context.Context,
	signupOrgUserReq db.SignupOrgUserReq,
) error {
	query := `
UPDATE
    org_users
SET
    name = $1,
    password_hash = $2,
    org_user_state = $3 -- ACTIVE_ORG_USER
WHERE
    id = (
        SELECT
            id
        FROM
            org_user_invites
        WHERE
            token = $4
    )
    AND org_user_state = $5 -- ADDED_ORG_USER
RETURNING
    id
`

	var orgUserID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		signupOrgUserReq.Name,
		signupOrgUserReq.PasswordHash,
		employer.ActiveOrgUserState,
		signupOrgUserReq.InviteToken,
		employer.AddedOrgUserState,
	).Scan(&orgUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrInviteTokenNotFound
		}

		p.log.Err("failed to signup org user", "error", err)
		return err
	}

	p.log.Dbg("org user signed up", "org_user_id", orgUserID)
	return nil
}

func (p *PG) UpdateOrgUser(
	ctx context.Context,
	updateOrgUserReq employer.UpdateOrgUserRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
WITH target_user AS (
    SELECT id, org_user_roles
    FROM org_users
    WHERE email = $1
      AND employer_id = $4::UUID
),
is_last_admin AS (
    SELECT 1
    FROM org_users
    WHERE employer_id = $4::UUID
      AND id != (SELECT id FROM target_user)
      AND 'ADMIN' = ANY(org_user_roles::org_user_roles[])
      AND org_user_state = $5 -- ACTIVE_ORG_USER
    LIMIT 1
),
updated_user AS (
    UPDATE org_users
    SET name = $2,
        org_user_roles = $3::org_user_roles[]
    WHERE id = (SELECT id FROM target_user)
      AND ('ADMIN' = ANY($3::org_user_roles[]) OR EXISTS (SELECT 1 FROM is_last_admin))
    RETURNING id
)
SELECT 
    CASE 
        WHEN (SELECT id FROM target_user) IS NULL THEN $6
        WHEN NOT EXISTS (SELECT 1 FROM updated_user) THEN $7
        ELSE (SELECT id FROM updated_user)::TEXT
    END AS result;
`

	const (
		userNotFound    = "USER_NOT_FOUND"
		lastActiveAdmin = "LAST_ACTIVE_ADMIN"
	)

	var orgUserIDStr string
	err := p.pool.QueryRow(
		ctx,
		query,
		updateOrgUserReq.Email,
		updateOrgUserReq.Name,
		updateOrgUserReq.Roles.StringArray(),
		orgUser.EmployerID,
		employer.ActiveOrgUserState,
		userNotFound,
		lastActiveAdmin,
	).Scan(&orgUserIDStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to update org user", "error", err)
		return uuid.UUID{}, err
	}

	if orgUserIDStr == userNotFound {
		return uuid.UUID{}, db.ErrNoOrgUser
	} else if orgUserIDStr == lastActiveAdmin {
		return uuid.UUID{}, db.ErrLastActiveAdmin
	}

	orgUserID, err := uuid.Parse(orgUserIDStr)
	if err != nil {
		p.log.Err("failed to parse org user id", "error", err)
		return uuid.UUID{}, err
	}

	return orgUserID, nil
}
