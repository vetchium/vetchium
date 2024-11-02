package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *PG) GetDomainNames(
	ctx context.Context,
	employerID uuid.UUID,
) ([]string, error) {
	query := "SELECT domain_name FROM domains WHERE employer_id = $1::UUID"

	rows, err := p.pool.Query(ctx, query, employerID)
	if err != nil {
		p.log.Error("failed to query domains", "error", err)
		return nil, err
	}

	domains, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		p.log.Error("failed to collect rows", "error", err)
		return nil, err
	}

	return domains, nil
}
