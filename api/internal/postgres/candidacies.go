package postgres

import (
	"context"
	"fmt"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetCandidaciesInfo(
	ctx context.Context,
	getCandidaciesInfoReq employer.GetCandidaciesInfoRequest,
) ([]employer.Candidacy, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []employer.Candidacy{}, db.ErrInternal
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
	i := 2

	if getCandidaciesInfoReq.RecruiterEmail != nil {
		query += fmt.Sprintf(` AND ou.email = $%d`, i)
		args = append(args, *getCandidaciesInfoReq.RecruiterEmail)
		i++
	}

	if getCandidaciesInfoReq.State != nil {
		query += fmt.Sprintf(` AND c.candidacy_state = $%d`, i)
		args = append(args, *getCandidaciesInfoReq.State)
		i++
	}

	if getCandidaciesInfoReq.PaginationKey != nil {
		query += fmt.Sprintf(` AND c.id > $%d`, i)
		args = append(args, *getCandidaciesInfoReq.PaginationKey)
		i++
	}

	query += " ORDER BY c.created_at DESC, c.id ASC "

	query += fmt.Sprintf(" LIMIT %d", getCandidaciesInfoReq.Limit)

	p.log.Dbg("query", "query", query)

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		p.log.Err("failed to query candidacies", "error", err)
		return []employer.Candidacy{}, db.ErrInternal
	}

	defer rows.Close()

	candidacies := []employer.Candidacy{}
	for rows.Next() {
		var candidacy employer.Candidacy
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
			return []employer.Candidacy{}, db.ErrInternal
		}

		candidacies = append(candidacies, candidacy)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over candidacies", "error", err)
		return []employer.Candidacy{}, db.ErrInternal
	}

	p.log.Dbg("candidacies", "candidacies", candidacies)

	return candidacies, nil
}

func (p *PG) GetMyCandidacies(
	ctx context.Context,
	getMyCandidaciesReq hub.MyCandidaciesRequest,
) ([]hub.MyCandidacy, error) {
	return []hub.MyCandidacy{}, nil
}
