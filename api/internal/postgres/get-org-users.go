package postgres

import (
	"context"
	"errors"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"

	"github.com/jackc/pgx/v5"
)

func (p *PG) GetOrgUserByEmail(
	ctx context.Context,
	email string,
) (db.OrgUserTO, error) {
	seekingOrgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("org user not found in context")
		return db.OrgUserTO{}, errors.New("org user not found in context")
	}

	var roles []string
	var orgUser db.OrgUserTO
	err := p.pool.QueryRow(ctx, `
SELECT id, name, email, password_hash, employer_id, org_user_roles, org_user_state, created_at
FROM org_users
WHERE email = $1
AND employer_id = $2
	`, email, seekingOrgUser.EmployerID).Scan(
		&orgUser.ID,
		&orgUser.Name,
		&orgUser.Email,
		&orgUser.PasswordHash,
		&orgUser.EmployerID,
		&roles,
		&orgUser.OrgUserState,
		&orgUser.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no org users found", "email", email)
			return db.OrgUserTO{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to scan org user", "error", err)
		return db.OrgUserTO{}, db.ErrInternal
	}

	orgUser.OrgUserRoles, err = p.convertToOrgUserRoles(roles)
	if err != nil {
		p.log.Err("failed to convert to org user roles", "error", err)
		return db.OrgUserTO{}, db.ErrInternal
	}

	return orgUser, nil
}
