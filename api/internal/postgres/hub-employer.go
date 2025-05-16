package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
)

// GetHubEmployerDetailsByDomain retrieves employer details based on a domain name.
// It returns db.ErrNoDomain if the domain is not found or not properly associated with an employer.
func (p *PG) GetHubEmployerDetailsByDomain(
	ctx context.Context,
	domainName string,
) (db.EmployerDetailsForHub, error) {
	var details db.EmployerDetailsForHub

	query := `
WITH domain_lookup AS (
    -- Find the domain and its associated employer_id
    SELECT id as domain_id, employer_id
    FROM domains
    WHERE domain_name = $1
),
employer_data AS (
    -- Get employer details if the domain is linked to an employer
    SELECT
        e.id AS employer_id,
        e.company_name
    FROM employers e
    JOIN domain_lookup dl ON e.id = dl.employer_id
    WHERE dl.employer_id IS NOT NULL -- Ensure the domain is actually linked to an employer
),
primary_domain_info AS (
    -- Find the primary domain for the identified employer
    SELECT
        ed.employer_id,
        pd.domain_name AS primary_domain_name
    FROM employer_data ed
    LEFT JOIN employer_primary_domains epd ON ed.employer_id = epd.employer_id
    LEFT JOIN domains pd ON epd.domain_id = pd.id -- Get the actual domain name
)
SELECT
    ed.employer_id,
    ed.company_name,
    COALESCE(pdi.primary_domain_name, '') AS primary_domain, -- Use empty string if no primary domain is set
    COALESCE(
        ARRAY_AGG(all_other_domains.domain_name) FILTER (WHERE all_other_domains.domain_name IS NOT NULL AND all_other_domains.domain_name != pdi.primary_domain_name),
        '{}'
    ) AS other_domains
FROM employer_data ed
LEFT JOIN primary_domain_info pdi ON ed.employer_id = pdi.employer_id
-- Join with all domains of this employer to list them as 'other' if they are not the primary
LEFT JOIN domains all_other_domains ON ed.employer_id = all_other_domains.employer_id
WHERE ed.employer_id IS NOT NULL -- This condition is effectively applied in employer_data CTE
GROUP BY ed.employer_id, ed.company_name, pdi.primary_domain_name
`
	err := p.pool.QueryRow(ctx, query, domainName).Scan(
		&details.EmployerID,
		&details.Name,
		&details.PrimaryDomain,
		&details.OtherDomains,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// This error means the query returned no rows, which implies:
			// 1. The domainName was not found in the 'domains' table.
			// 2. The domain was found, but it's not associated with any employer_id.
			// 3. The employer_id associated with the domain does not exist in the 'employers' table (FK violation, unlikely).
			// In any of these cases, it's appropriate to return db.ErrNoDomain.
			return db.EmployerDetailsForHub{}, db.ErrNoDomain
		}
		// For any other database error, log it and return a generic internal error.
		p.log.Err("db error", "error", err, "domain", domainName)
		return db.EmployerDetailsForHub{}, db.ErrInternal
	}

	p.log.Dbg("got employer details", "domain", domainName, "details", details)
	return details, nil
}
