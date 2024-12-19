package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetAssessment(
	ctx context.Context,
	req employer.GetAssessmentRequest,
) (employer.Assessment, error) {
	p.log.Dbg("GetAssessment", "req", req)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("error getting org user")
		return employer.Assessment{}, db.ErrInternal
	}

	query := `
SELECT
	id,
	positives,
	negatives,
	overall_assessment,
	feedback_to_candidate,
	interviewers_decision,
	feedback_submitted_by,
	feedback_submitted_at
FROM interviews
WHERE id = $1
AND employer_id = $2
`

	p.log.Dbg("query", "query", query)

	var assessment employer.Assessment
	err := p.pool.QueryRow(ctx, query, req.InterviewID, orgUser.EmployerID).
		Scan(
			&assessment.InterviewID,
			&assessment.Decision,
			&assessment.Positives,
			&assessment.Negatives,
			&assessment.OverallAssessment,
			&assessment.FeedbackToCandidate,
			&assessment.FeedbackSubmittedBy,
			&assessment.FeedbackSubmittedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no interview found", "interview_id", req.InterviewID)
			return employer.Assessment{}, db.ErrNoInterview
		}

		p.log.Err("error getting assessment", "error", err)
		return employer.Assessment{}, db.ErrInternal
	}

	return assessment, nil
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

	const (
		noInterview    = "NO_INTERVIEW"
		notInterviewer = "NOT_INTERVIEWER"
		wrongState     = "WRONG_STATE"
		success        = "SUCCESS"
	)

	query := `
WITH validation AS (
	SELECT
		i.id,
		i.interview_state,
		EXISTS (
			SELECT 1 FROM interview_interviewers ii
			WHERE ii.interview_id = i.id
			AND ii.interviewer_id = $7
		) as is_interviewer
	FROM interviews i
	WHERE i.id = $1
)
UPDATE interviews i
SET
	positives = $2,
	negatives = $3,
	overall_assessment = $4,
	feedback_to_candidate = $5,
	interviewers_decision = $6,
	feedback_submitted_by = $7,
	feedback_submitted_at = NOW(),
	updated_at = NOW()
FROM validation v
WHERE i.id = v.id
AND v.is_interviewer = true
AND v.interview_state = $8
RETURNING
	CASE
		WHEN v.id IS NULL THEN '` + noInterview + `'
		WHEN v.is_interviewer = false THEN '` + notInterviewer + `'
		WHEN v.interview_state != $8 THEN '` + wrongState + `'
		ELSE '` + success + `'
	END as result
`

	p.log.Dbg("query", "query", query)

	var result string
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
		string(common.ScheduledInterviewState),
	).Scan(&result)
	if err != nil {
		p.log.Err("error putting assessment", "error", err)
		return db.ErrInternal
	}

	switch result {
	case success:
		return nil
	case noInterview:
		return db.ErrNoInterview
	case notInterviewer:
		return db.ErrNotAnInterviewer
	case wrongState:
		return db.ErrStateMismatch
	default:
		return db.ErrInternal
	}
}
