package postgres

import (
	"context"
	"errors"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"

	"github.com/jackc/pgx/v5"
)

func (p *PG) GetOrgUsersByEmails(
	ctx context.Context,
	emails []string,
) ([]db.OrgUserTO, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("org user not found in context")
		return []db.OrgUserTO{}, errors.New("org user not found in context")
	}

	rows, err := p.pool.Query(ctx, `
SELECT id, name, email, password_hash, employer_id, org_user_roles, org_user_state, created_at
FROM org_users
WHERE email = ANY($1)
AND employer_id = $2
	`, emails, orgUser.EmployerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no org users found", "emails", emails)
			return []db.OrgUserTO{}, db.ErrNoOrgUser
		}

		p.log.Err("failed to get org users by emails", "error", err)
		return []db.OrgUserTO{}, err
	}

	orgUsers := make([]db.OrgUserTO, 0, len(emails))
	for rows.Next() {
		var orgUser db.OrgUserTO
		if err := rows.Scan(
			&orgUser.ID,
			&orgUser.Name,
			&orgUser.Email,
			&orgUser.PasswordHash,
			&orgUser.EmployerID,
			&orgUser.OrgUserRoles,
			&orgUser.OrgUserState,
			&orgUser.CreatedAt,
		); err != nil {
			p.log.Err("failed to scan org user", "error", err)
			return []db.OrgUserTO{}, err
		}

		orgUsers = append(orgUsers, orgUser)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over org users", "error", err)
		return []db.OrgUserTO{}, err
	}

	return orgUsers, nil
}
