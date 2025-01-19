package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) FilterCandidacyInfos(
	ctx context.Context,
	filterCandidacyInfosReq employer.FilterCandidacyInfosRequest,
) ([]common.Candidacy, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []common.Candidacy{}, db.ErrInternal
	}

	var args []interface{}

	query := `
SELECT c.id, o.id, o.title, o.jd, c.candidacy_state, hu.full_name, hu.handle
FROM candidacies c
JOIN openings o ON c.employer_id = o.employer_id AND c.opening_id = o.id 
JOIN applications a ON c.application_id = a.id
JOIN hub_users hu ON a.hub_user_id = hu.id
JOIN org_users ou ON o.recruiter = ou.id
WHERE c.employer_id = $1
	`
	args = append(args, orgUser.EmployerID)

	if filterCandidacyInfosReq.OpeningID != nil {
		query += fmt.Sprintf(` AND c.opening_id = $%d`, len(args)+1)
		args = append(args, *filterCandidacyInfosReq.OpeningID)
	}

	if filterCandidacyInfosReq.RecruiterEmail != nil {
		query += fmt.Sprintf(` AND ou.email = $%d`, len(args)+1)
		args = append(args, *filterCandidacyInfosReq.RecruiterEmail)
	}

	if filterCandidacyInfosReq.State != nil {
		query += fmt.Sprintf(` AND c.candidacy_state = $%d`, len(args)+1)
		args = append(args, *filterCandidacyInfosReq.State)
	}

	if filterCandidacyInfosReq.PaginationKey != nil {
		query += fmt.Sprintf(` AND c.id > $%d`, len(args)+1)
		args = append(args, *filterCandidacyInfosReq.PaginationKey)
	}

	query += " ORDER BY c.created_at DESC, c.id ASC "

	query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
	args = append(args, filterCandidacyInfosReq.Limit)

	p.log.Dbg("query", "query", query)

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to query candidacies", "error", err)
		return []common.Candidacy{}, db.ErrInternal
	}

	defer rows.Close()

	candidacies := []common.Candidacy{}
	for rows.Next() {
		var candidacy common.Candidacy
		err := rows.Scan(
			&candidacy.CandidacyID,
			&candidacy.OpeningID,
			&candidacy.OpeningTitle,
			&candidacy.OpeningDescription,
			&candidacy.CandidacyState,
			&candidacy.ApplicantName,
			&candidacy.ApplicantHandle,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy", "error", err)
			return []common.Candidacy{}, db.ErrInternal
		}

		candidacies = append(candidacies, candidacy)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over candidacies", "error", err)
		return []common.Candidacy{}, db.ErrInternal
	}

	p.log.Dbg("candidacies", "candidacies", candidacies)

	return candidacies, nil
}

func (p *PG) GetMyCandidacies(
	ctx context.Context,
	getMyCandidaciesReq hub.MyCandidaciesRequest,
) ([]hub.MyCandidacy, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return []hub.MyCandidacy{}, db.ErrInternal
	}

	query := `
SELECT c.id, e.company_name, d.domain_name, o.id, o.title, o.jd, c.candidacy_state
FROM candidacies c
JOIN openings o ON c.opening_id = o.id
JOIN employers e ON o.employer_id = e.id
JOIN domains d ON e.id = d.employer_id
JOIN employer_primary_domains epd ON d.id = epd.domain_id AND epd.employer_id = e.id
JOIN applications a ON c.application_id = a.id
WHERE a.hub_user_id = $1
`

	var args []interface{}
	args = append(args, hubUser.ID)

	if getMyCandidaciesReq.CandidacyStates != nil {
		stateParams := make(
			[]string,
			0,
			len(getMyCandidaciesReq.CandidacyStates),
		)
		for j, state := range getMyCandidaciesReq.CandidacyStates {
			paramNum := len(args) + 1 + j
			stateParams = append(
				stateParams,
				fmt.Sprintf("$%d::candidacy_states", paramNum),
			)
			args = append(args, string(state))
		}
		query += fmt.Sprintf(
			" AND c.candidacy_state = ANY(ARRAY[%s])",
			strings.Join(stateParams, ","),
		)
	}

	if getMyCandidaciesReq.PaginationKey != nil {
		query += fmt.Sprintf(` AND c.id > $%d`, len(args)+1)
		args = append(args, *getMyCandidaciesReq.PaginationKey)
	}

	query += " ORDER BY c.created_at DESC, c.id ASC "

	query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
	args = append(args, getMyCandidaciesReq.Limit)

	p.log.Dbg("query", "query", query, "args", args)

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to query candidacies", "error", err)
		return []hub.MyCandidacy{}, db.ErrInternal
	}
	defer rows.Close()

	candidacies := []hub.MyCandidacy{}
	for rows.Next() {
		var candidacy hub.MyCandidacy
		err := rows.Scan(
			&candidacy.CandidacyID,
			&candidacy.CompanyName,
			&candidacy.CompanyDomain,
			&candidacy.OpeningID,
			&candidacy.OpeningTitle,
			&candidacy.OpeningDescription,
			&candidacy.CandidacyState,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy", "error", err)
			return []hub.MyCandidacy{}, db.ErrInternal
		}

		candidacies = append(candidacies, candidacy)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over candidacies", "error", err)
		return []hub.MyCandidacy{}, db.ErrInternal
	}

	p.log.Dbg("candidacies", "candidacies", candidacies)

	return candidacies, nil
}
