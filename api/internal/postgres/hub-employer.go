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

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hubUserID from context",
			"error", err,
			"domain", domainName)
		return db.EmployerDetailsForHub{}, db.ErrInternal
	}

	query := `
WITH domain_lookup AS (
    SELECT id as domain_id, employer_id
    FROM domains
    WHERE domain_name = $1
),
employer_data AS (
    SELECT
        e.id AS employer_id,
        e.company_name,
        e.employer_state
    FROM employers e
    JOIN domain_lookup dl ON e.id = dl.employer_id
    WHERE dl.employer_id IS NOT NULL
),
primary_domain_info AS (
    SELECT
        ed.employer_id,
        pd.domain_name AS primary_domain_name
    FROM employer_data ed
    LEFT JOIN employer_primary_domains epd ON ed.employer_id = epd.employer_id
    LEFT JOIN domains pd ON epd.domain_id = pd.id
)
SELECT
    ed.employer_id,
    ed.company_name,
    COALESCE(pdi.primary_domain_name, '') AS primary_domain,
    COALESCE(
        ARRAY_AGG(all_other_domains.domain_name) FILTER (WHERE all_other_domains.domain_name IS NOT NULL AND all_other_domains.domain_name != pdi.primary_domain_name),
        '{}'
    ) AS other_domains,
    (ed.employer_state = 'ONBOARDED') AS is_onboarded,
    EXISTS (
        SELECT 1
        FROM org_following_relationships ofr
        WHERE ofr.employer_id = ed.employer_id
        AND ofr.hub_user_id = $2::uuid
    ) AS is_following
FROM employer_data ed
LEFT JOIN primary_domain_info pdi ON ed.employer_id = pdi.employer_id
LEFT JOIN domains all_other_domains ON ed.employer_id = all_other_domains.employer_id
WHERE ed.employer_id IS NOT NULL
GROUP BY ed.employer_id, ed.company_name, pdi.primary_domain_name, ed.employer_state
`

	err = p.pool.QueryRow(ctx, query, domainName, hubUserID).Scan(
		&details.EmployerID,
		&details.Name,
		&details.PrimaryDomain,
		&details.OtherDomains,
		&details.IsOnboarded,
		&details.IsFollowing,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.EmployerDetailsForHub{}, db.ErrNoDomain
		}
		p.log.Err("db error querying employer details by domain",
			"error", err, "domain", domainName)
		return db.EmployerDetailsForHub{}, db.ErrInternal
	}

	p.log.Dbg("got employer details", "domain", domainName, "details", details)
	return details, nil
}
