package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) AddOrgUser(
	ctx context.Context,
	req db.AddOrgUserReq,
) (uuid.UUID, error) {
	query := `
	INSERT INTO org_users (email, employer_id, org_user_roles, org_user_state)
		VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	var id uuid.UUID
	err := p.pool.QueryRow(ctx, query, req.Email, req.EmployerID, req.OrgUserRoles, req.OrgUserState).
		Scan(&id)
	if err != nil {
		p.log.Error("failed to add org user", "error", err)
		return uuid.Nil, err
	}

	return id, nil
}

func (p *PG) DisableOrgUser(
	ctx context.Context,
	req db.DisableOrgUserReq,
) error {
	const (
		userNotFound    = "USER_NOT_FOUND"
		lastActiveAdmin = "LAST_ACTIVE_ADMIN"
	)
	query := `
WITH target_user AS (
    SELECT id
    FROM org_users
    WHERE employer_id = $1
      AND email = $2
),
is_last_admin AS (
    SELECT 1
    FROM org_users
    WHERE employer_id = $1
      AND id != (SELECT id FROM target_user)
      AND 'ADMIN' = ANY(org_user_roles)
      AND org_user_state = $3 -- ACTIVE_ORG_USER
    LIMIT 1
),
updated_user AS (
    UPDATE org_users
    SET org_user_state = $4 -- DISABLED_ORG_USER
    WHERE id = (SELECT id FROM target_user)
      AND org_user_state != $4 -- DISABLED_ORG_USER
      AND NOT EXISTS (SELECT 1 FROM is_last_admin)
    RETURNING id
),
deleted_tokens AS (
    DELETE FROM org_user_tokens
    WHERE org_user_id = (SELECT id FROM updated_user)
)
SELECT 
    CASE 
        WHEN (SELECT id FROM target_user) IS NULL THEN $5
        WHEN NOT EXISTS (SELECT 1 FROM updated_user) THEN $6
        ELSE (SELECT id FROM updated_user)::TEXT
    END AS result;
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.EmployerID,
		req.Email,
		db.ActiveOrgUserState,
		db.DisabledOrgUserState,
		userNotFound,
		lastActiveAdmin,
	).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ErrNoOrgUser
		}

		p.log.Error("failed to disable org user", "error", err)
		return err
	}

	switch result {
	case userNotFound:
		return db.ErrNoOrgUser
	case lastActiveAdmin:
		return db.ErrLastActiveAdmin
	default:
		orgUserID, err := uuid.Parse(result)
		if err != nil {
			p.log.Error("failed to parse org user id", "error", err)
			return err
		}

		p.log.Debug("org user disabled", "org_user_id", orgUserID)
		return nil
	}
}
