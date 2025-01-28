package postgres

import (
	"context"
	"database/sql"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) OfferToCandidate(
	ctx context.Context,
	request employer.OfferToCandidateRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	query := `
UPDATE candidacy
SET status = $1
WHERE id = $2
AND employer_id = $3
AND status = $4
RETURNING id
`

	err = tx.QueryRow(
		ctx,
		query,
		common.OfferedCandidacyState,
		request.CandidacyID,
		orgUser.EmployerID,
		common.InterviewingCandidacyState,
	).Scan(&request.CandidacyID)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Dbg("not found", "candidacy_id", request.CandidacyID)
			return db.ErrNoCandidacy
		}
		p.log.Err("failed to update candidacy status", "error", err)
		return db.ErrInternal
	}

	interviewUpdateQuery := `
UPDATE interviews
SET interview_state = $1
WHERE candidacy_id = $2
AND interview_state = $3
`
	_, err = tx.Exec(
		ctx,
		interviewUpdateQuery,
		common.CancelledInterviewState,
		request.CandidacyID,
		common.ScheduledInterviewState,
	)
	if err != nil {
		p.log.Err("failed to update interview state", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
