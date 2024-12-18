package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetAssessment(
	ctx context.Context,
	req employer.GetAssessmentRequest,
) (employer.Assessment, error) {
	return employer.Assessment{}, nil
}

func (p *PG) PutAssessment(
	ctx context.Context,
	req employer.Assessment,
) error {
	p.log.Dbg("PutAssessment", "req", req)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("error getting org user")
		return db.ErrInternal
	}

	// TODO: Add checks for interviewer, interview state, etc.
	query := `
UPDATE interviews
SET
	positives = $2,
	negatives = $3,
	overall_assessment = $4,
	feedback_to_candidate = $5,
	decision = $6,
WHERE id = $1
	`

	var status string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.InterviewID,
		req.Positives,
		req.Negatives,
		req.OverallAssessment,
		req.FeedbackToCandidate,
		req.Decision,
		orgUser.ID,
	).Scan(&status)
	if err != nil {
		p.log.Err("error putting assessment", "error", err)
		return err
	}

	return nil
}
