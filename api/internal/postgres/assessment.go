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
	interviewers_decision,
	positives,
	negatives,
	overall_assessment,
	feedback_to_candidate,
	feedback_submitted_by,
	feedback_submitted_at
FROM interviews
WHERE id = $1
AND employer_id = $2
`

	p.log.Dbg("query", "query", query)

	var assessment employer.Assessment
	var decision *common.InterviewersDecision
	var positives, negatives, overallAssessment, feedbackToCandidate *string
	err := p.pool.QueryRow(ctx, query, req.InterviewID, orgUser.EmployerID).
		Scan(
			&assessment.InterviewID,
			&decision,
			&positives,
			&negatives,
			&overallAssessment,
			&feedbackToCandidate,
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

	// Convert nullable fields to empty values if NULL
	if decision != nil {
		assessment.Decision = *decision
	}
	if positives != nil {
		assessment.Positives = *positives
	}
	if negatives != nil {
		assessment.Negatives = *negatives
	}
	if overallAssessment != nil {
		assessment.OverallAssessment = *overallAssessment
	}
	if feedbackToCandidate != nil {
		assessment.FeedbackToCandidate = *feedbackToCandidate
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

	p.log.Dbg(
		"PutAssessment context",
		"org_user_id",
		orgUser.ID,
		"email",
		orgUser.Email,
	)

	const (
		noInterview    = "NO_INTERVIEW"
		notInterviewer = "NOT_INTERVIEWER"
		wrongState     = "WRONG_STATE"
		success        = "SUCCESS"
	)

	query := `
WITH interviewer_debug AS (
	SELECT ii.interviewer_id, ou.id as org_user_id, ou.email
	FROM interview_interviewers ii
	JOIN org_users ou ON ou.id = ii.interviewer_id
	WHERE ii.interview_id = $1
),
validation AS (
	SELECT
		i.id,
		i.interview_state,
		i.employer_id,
		EXISTS (
			SELECT 1 FROM interview_interviewers ii
			WHERE ii.interview_id = i.id
			AND ii.interviewer_id = $7
		) as is_interviewer,
		(SELECT json_agg(row_to_json(d)) FROM interviewer_debug d) as debug_interviewers
	FROM interviews i
	WHERE i.id = $1
	AND i.employer_id = $9
),
validation_result AS (
	SELECT
		id,
		interview_state,
		is_interviewer,
		debug_interviewers,
		CASE
			WHEN id IS NULL THEN '` + noInterview + `'
			WHEN is_interviewer = false THEN '` + notInterviewer + `'
			WHEN interview_state != $8 THEN '` + wrongState + `'
			ELSE '` + success + `'
		END as validation_status
	FROM validation
),
update_result AS (
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
	FROM validation_result v
	WHERE i.id = v.id
	AND v.validation_status = '` + success + `'
	RETURNING i.id
)
SELECT
	v.validation_status as result,
	v.interview_state as debug_state,
	v.is_interviewer as debug_is_interviewer,
	$8 as debug_expected_state,
	v.debug_interviewers as debug_interviewers
FROM validation_result v
LIMIT 1
`

	p.log.Dbg("query", "query", query)

	var (
		result             string
		debugState         string
		debugIsInterviewer bool
		debugExpectedState string
		debugInterviewers  []byte
	)
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
		orgUser.EmployerID,
	).Scan(&result, &debugState, &debugIsInterviewer, &debugExpectedState, &debugInterviewers)
	if err != nil {
		p.log.Err("error putting assessment", "error", err)
		return db.ErrInternal
	}

	p.log.Dbg("assessment validation result",
		"result", result,
		"state", debugState,
		"is_interviewer", debugIsInterviewer,
		"expected_state", debugExpectedState,
		"debug_interviewers", string(debugInterviewers))

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
