package postgres

import (
	"context"
	"fmt"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
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
SELECT c.id, o.id, o.title, o.description, c.candidacy_state, a.name, a.handle
FROM candidacies c
JOIN openings o ON c.opening_id = o.id
JOIN applicants a ON c.applicant_id = a.id
WHERE c.employer_id = $1
	`
	args = append(args, orgUser.EmployerID)
	i := 2

	if getCandidaciesInfoReq.RecruiterID != nil {
		query += fmt.Sprintf(` AND recruiter_id = $%d`, i)
		args = append(args, *getCandidaciesInfoReq.RecruiterID)
		i++
	}

	if getCandidaciesInfoReq.State != nil {
		query += fmt.Sprintf(` AND candidacy_state = $%d`, i)
		args = append(args, *getCandidaciesInfoReq.State)
		i++
	}

	if getCandidaciesInfoReq.PaginationKey != nil {
		query += fmt.Sprintf(` AND id > $%d`, i)
		args = append(args, *getCandidaciesInfoReq.PaginationKey)
		i++
	}

	query += " ORDER BY created_at DESC, id DESC "
	query += fmt.Sprintf("LIMIT %d", getCandidaciesInfoReq.Limit)

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
