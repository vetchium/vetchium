package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) AddInterview(
	ctx context.Context,
	req employer.AddInterviewRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	const (
		statusNoCandidacy = "no_candidacy"
		statusWrongState  = "wrong_state"
		statusOK          = "ok"
	)

	query := `
WITH state_validation AS (
	SELECT
		CASE
			WHEN NOT EXISTS (
				SELECT 1 FROM candidacies c
				WHERE c.id = $1
				AND c.employer_id = $2
			) THEN $9
			WHEN EXISTS (
				SELECT 1 FROM candidacies c
				JOIN applications a ON c.application_id = a.id
				JOIN openings o ON c.opening_id = o.id
				WHERE c.id = $1
				AND c.employer_id = $2
				AND (
					c.candidacy_state != $10::candidacy_states
					OR a.application_state != $11::application_states
					OR o.state != $12::opening_states
				)
			) THEN $13
			ELSE $14
		END as status
),
interview_insert AS (
	INSERT INTO interviews (
		candidacy_id,
		interview_type,
		interview_state,
		start_time,
		end_time,
		description,
		created_by,
		employer_id,
		candidate_rsvp
	)
	SELECT
		$1,                             -- candidacy_id
		$3::interview_types,            -- interview_type
		$4::interview_states,           -- interview_state
		$5::timestamptz,               -- start_time
		$6::timestamptz,               -- end_time
		$7,                            -- description
		$8,                            -- created_by
		$2,                            -- employer_id
		'NOT_SET'::rsvp_status         -- candidate_rsvp default
	WHERE (SELECT status FROM state_validation) = $14
	RETURNING id
)
SELECT
	COALESCE(
		(SELECT id::text FROM interview_insert),
		(SELECT status FROM state_validation)
	)
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.CandidacyID,                   // $1
		orgUser.EmployerID,                // $2
		req.InterviewType,                 // $3
		common.ScheduledInterviewState,    // $4
		req.StartTime,                     // $5
		req.EndTime,                       // $6
		req.Description,                   // $7
		orgUser.ID,                        // $8
		statusNoCandidacy,                 // $9
		common.InterviewingCandidacyState, // $10
		common.AppliedAppState,            // $11
		common.ActiveOpening,              // $12
		statusWrongState,                  // $13
		statusOK,                          // $14
	).Scan(&result)
	if err != nil {
		p.log.Err("failed to add interview", "error", err)
		return uuid.UUID{}, db.ErrInternal
	}

	switch result {
	case statusNoCandidacy:
		return uuid.UUID{}, db.ErrNoCandidacy
	case statusWrongState:
		return uuid.UUID{}, db.ErrInvalidCandidacyState
	default:
		interviewID, err := uuid.Parse(result)
		if err != nil {
			p.log.Err("failed to parse interview id", "error", err)
			return uuid.UUID{}, db.ErrInternal
		}
		return interviewID, nil
	}
}

func (p *PG) GetInterviewsByOpening(
	ctx context.Context,
	req employer.GetInterviewsByOpeningRequest,
) ([]employer.Interview, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	// First verify the opening belongs to the employer
	verifyQuery := `
	SELECT EXISTS (
		SELECT 1 FROM openings 
		WHERE employer_id = $1 
		AND id = $2
	)`

	var exists bool
	err := p.pool.QueryRow(ctx, verifyQuery, orgUser.EmployerID, req.OpeningID).
		Scan(&exists)
	if err != nil {
		p.log.Err("failed to verify opening ownership", "error", err)
		return nil, db.ErrInternal
	}
	if !exists {
		return nil, db.ErrNoOpening
	}

	query := `
	WITH interview_data AS (
		SELECT
			i.id,
			i.interview_type,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.description,
			i.interviewers_decision,
			i.interviewers_assessment as positives,
			i.interviewers_assessment as negatives,
			i.interviewers_assessment as overall_assessment,
			i.feedback_to_candidate,
			i.created_at,
			i.feedback_submitted_at,
			CASE 
				WHEN fb.id IS NOT NULL THEN
					json_build_object(
						'name', fb.name,
						'email', fb.email
					)
				ELSE NULL
			END as feedback_submitted_by,
			COALESCE(
				json_agg(
					json_build_object(
						'name', ou.name,
						'email', ou.email
					) 
					ORDER BY ou.email
				) FILTER (WHERE ou.id IS NOT NULL),
				'[]'::json
			) as interviewers
		FROM interviews i
		JOIN candidacies c ON i.candidacy_id = c.id
		LEFT JOIN interview_interviewers ii ON i.id = ii.interview_id
		LEFT JOIN org_users ou ON ii.interviewer_id = ou.id
		LEFT JOIN org_users fb ON i.feedback_submitted_by = fb.id
		WHERE c.opening_id = $1
		AND i.employer_id = $2
		AND ($3::uuid IS NULL OR i.id::text > $3)
		GROUP BY 
			i.id,
			i.interview_type,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.description,
			i.interviewers_decision,
			i.interviewers_assessment,
			i.feedback_to_candidate,
			i.created_at,
			fb.id,
			fb.name,
			fb.email,
			i.feedback_submitted_at
		ORDER BY i.start_time ASC, i.id ASC
		LIMIT $4
	)
	SELECT 
		id,
		interview_type,
		interview_state,
		start_time,
		end_time,
		description,
		interviewers,
		interviewers_decision,
		positives,
		negatives,
		overall_assessment,
		feedback_to_candidate,
		created_at,
		feedback_submitted_by,
		feedback_submitted_at
	FROM interview_data
	`

	var paginationKey *uuid.UUID
	if req.PaginationKey != "" {
		parsed, err := uuid.Parse(req.PaginationKey)
		if err != nil {
			p.log.Err("failed to parse pagination key", "error", err)
			return nil, db.ErrInvalidPaginationKey
		}
		paginationKey = &parsed
	}

	rows, err := p.pool.Query(ctx, query,
		req.OpeningID,      // $1
		orgUser.EmployerID, // $2
		paginationKey,      // $3
		req.Limit,          // $4
	)
	if err != nil {
		p.log.Err("failed to get interviews by opening", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	interviews := []employer.Interview{}
	for rows.Next() {
		var interview employer.Interview
		var feedbackSubmittedBy *struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		err := rows.Scan(
			&interview.InterviewID,
			&interview.InterviewType,
			&interview.InterviewState,
			&interview.StartTime,
			&interview.EndTime,
			&interview.Description,
			&interview.Interviewers,
			&interview.InterviewersDecision,
			&interview.Positives,
			&interview.Negatives,
			&interview.OverallAssessment,
			&interview.FeedbackToCandidate,
			&interview.CreatedAt,
			&feedbackSubmittedBy,
			&interview.FeedbackSubmittedAt,
		)
		if err != nil {
			p.log.Err("failed to scan interview", "error", err)
			return nil, db.ErrInternal
		}

		if feedbackSubmittedBy != nil {
			interview.FeedbackSubmittedBy = &employer.OrgUserTiny{
				Name:  feedbackSubmittedBy.Name,
				Email: feedbackSubmittedBy.Email,
			}
		}

		interviews = append(interviews, interview)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over interviews", "error", err)
		return nil, db.ErrInternal
	}

	return interviews, nil
}

func (p *PG) GetInterviewsByCandidacy(
	ctx context.Context,
	req employer.GetInterviewsByCandidacyRequest,
) ([]employer.Interview, error) {
	return nil, nil
}
