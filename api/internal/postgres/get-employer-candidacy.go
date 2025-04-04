package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func (p *PG) GetEmployerCandidacyInfo(
	ctx context.Context,
	getCandidacyInfoReq common.GetCandidacyInfoRequest,
) (employer.Candidacy, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.Candidacy{}, db.ErrInternal
	}

	query := `
SELECT c.id, o.id, o.title, o.jd, c.candidacy_state, hu.full_name, hu.handle
FROM candidacies c
JOIN openings o ON c.opening_id = o.id
JOIN applications a ON c.application_id = a.id
JOIN hub_users hu ON a.hub_user_id = hu.id
WHERE c.id = $1
AND o.employer_id = $2
`

	var candidacy employer.Candidacy
	err := p.pool.QueryRow(
		ctx,
		query,
		getCandidacyInfoReq.CandidacyID,
		orgUser.EmployerID,
	).Scan(
		&candidacy.CandidacyID,
		&candidacy.OpeningID,
		&candidacy.OpeningTitle,
		&candidacy.OpeningDescription,
		&candidacy.CandidacyState,
		&candidacy.ApplicantName,
		&candidacy.ApplicantHandle,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg(
				"no candidacy found",
				"candidacy_id",
				getCandidacyInfoReq.CandidacyID,
				"employer_id",
				orgUser.EmployerID,
			)
			return employer.Candidacy{}, db.ErrNoCandidacy
		}
		p.log.Err("failed to scan candidacy", "error", err)
		return employer.Candidacy{}, db.ErrInternal
	}

	return candidacy, nil
}
