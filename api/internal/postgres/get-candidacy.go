package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) GetEmployerCandidacyInfo(
	ctx context.Context,
	getCandidacyInfoReq common.GetCandidacyInfoRequest,
) (common.Candidacy, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return common.Candidacy{}, db.ErrInternal
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

	var candidacy common.Candidacy
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
			return common.Candidacy{}, db.ErrNoCandidacy
		}
		p.log.Err("failed to scan candidacy", "error", err)
		return common.Candidacy{}, db.ErrInternal
	}

	return candidacy, nil
}

func (p *PG) GetHubCandidacyInfo(
	ctx context.Context,
	getCandidacyInfoReq common.GetCandidacyInfoRequest,
) (common.Candidacy, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return common.Candidacy{}, db.ErrInternal
	}

	query := `
SELECT c.id, o.id, o.title, o.jd, c.candidacy_state, hu.full_name, hu.handle
FROM candidacies c
JOIN openings o ON c.opening_id = o.id
JOIN applications a ON c.application_id = a.id
JOIN hub_users hu ON a.hub_user_id = hu.id
		WHERE c.id = $1
		AND a.hub_user_id = $2
	`

	var candidacy common.Candidacy
	err := p.pool.QueryRow(
		ctx,
		query,
		getCandidacyInfoReq.CandidacyID,
		hubUser.ID,
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
			return common.Candidacy{}, db.ErrNoCandidacy
		}
		p.log.Err("failed to scan candidacy", "error", err)
		return common.Candidacy{}, db.ErrInternal
	}
	return candidacy, nil
}
