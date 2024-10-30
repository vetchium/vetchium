package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/psankar/vetchi/api/pkg/vetchi"
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
