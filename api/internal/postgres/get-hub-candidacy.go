package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (p *PG) GetHubCandidacyInfo(
	ctx context.Context,
	getCandidacyInfoReq common.GetCandidacyInfoRequest,
) (hub.MyCandidacy, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return hub.MyCandidacy{}, db.ErrInternal
	}

	query := `
SELECT 
	c.id as candidacy_id,
	e.company_name,
	d.domain_name as company_domain,
	o.id as opening_id,
	o.title as opening_title,
	o.jd as opening_description,
	c.candidacy_state
FROM candidacies c
JOIN applications a ON a.id = c.application_id
JOIN openings o ON o.employer_id = c.employer_id AND o.id = c.opening_id
JOIN employers e ON e.id = c.employer_id
JOIN employer_primary_domains epd ON epd.employer_id = e.id
JOIN domains d ON d.id = epd.domain_id
WHERE c.id = $1 AND a.hub_user_id = $2
`

	var candidacy hub.MyCandidacy
	err := p.pool.QueryRow(
		ctx,
		query,
		getCandidacyInfoReq.CandidacyID,
		hubUser.ID,
	).Scan(
		&candidacy.CandidacyID,
		&candidacy.CompanyName,
		&candidacy.CompanyDomain,
		&candidacy.OpeningID,
		&candidacy.OpeningTitle,
		&candidacy.OpeningDescription,
		&candidacy.CandidacyState,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return hub.MyCandidacy{}, db.ErrNoCandidacy
		}
		p.log.Err("failed to scan candidacy", "error", err)
		return hub.MyCandidacy{}, db.ErrInternal
	}
	return candidacy, nil
}
