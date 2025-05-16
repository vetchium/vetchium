package postgres

import (
	"context"

	"github.com/vetchium/vetchium/typespec/common"
)

func (p *PG) GetEmployerActiveJobCount(
	ctx context.Context,
	domain string,
) (uint32, error) {
	query := `
SELECT COUNT(*) FROM openings
WHERE employer_id = (
	SELECT employer_id FROM domains WHERE domain_name = $1
) AND state = $2
`

	var count uint32
	err := p.pool.QueryRow(ctx, query, domain, common.ActiveOpening).
		Scan(&count)
	if err != nil {
		p.log.Err("failed to get employer active job count", "error", err)
		return 0, err
	}

	return count, nil
}

func (p *PG) GetEmployerEmployeeCount(
	ctx context.Context,
	domain string,
) (uint32, error) {
	query := `
SELECT COUNT(DISTINCT huo.hub_user_id)
FROM hub_users_official_emails huo
JOIN domains d ON huo.domain_id = d.id
WHERE d.employer_id = (
    SELECT employer_id FROM domains WHERE domain_name = $1
)
`

	var count uint32
	if err := p.pool.QueryRow(ctx, query, domain).Scan(&count); err != nil {
		p.log.Err("failed to get employer employee count", "error", err)
		return 0, err
	}

	return count, nil
}
