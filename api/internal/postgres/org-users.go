package postgres

import (
	"context"

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
