package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) FilterEmployers(
	ctx context.Context,
	req hub.FilterEmployersRequest,
) ([]hub.HubEmployer, error) {

	// TODO: This query will not work for non-ASCII domains
	query := `
WITH filtered_employers AS (
    SELECT DISTINCT ON (COALESCE(e.company_name, d.domain_name))
        COALESCE(e.company_name, d.domain_name) as name,
        d.domain_name as domain
    FROM employers e
    FULL OUTER JOIN domains d ON e.id = d.employer_id
    WHERE (e.employer_state = $1 OR e.employer_state IS NULL)
    AND (
        lower(COALESCE(e.company_name, d.domain_name)) LIKE lower($2) OR
        lower(d.domain_name) LIKE lower($2)
    )
    ORDER BY COALESCE(e.company_name, d.domain_name)
)
SELECT name, domain
FROM filtered_employers
LIMIT 8
`

	// Add % wildcards to the search term for LIKE pattern matching
	searchTerm := "%" + req.Prefix + "%"

	employers := make([]hub.HubEmployer, 0)
	rows, err := p.pool.Query(ctx, query, db.OnboardedEmployerState, searchTerm)
	if err != nil {
		p.log.Err("failed to execute filter employers query", "error", err)
		return employers, db.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var employer hub.HubEmployer
		err = rows.Scan(&employer.Name, &employer.Domain)
		if err != nil {
			p.log.Err("failed to scan employer row", "error", err)
			return employers, db.ErrInternal
		}
		employers = append(employers, employer)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("error occurred while iterating rows", "error", err)
		return employers, db.ErrInternal
	}

	return employers, nil
}
