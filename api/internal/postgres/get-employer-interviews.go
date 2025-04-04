package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/employer"
)

func (p *PG) GetEmployerInterviewsByOpening(
	ctx context.Context,
	req employer.GetEmployerInterviewsByOpeningRequest,
) ([]employer.EmployerInterview, error) {
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
			i.positives,
			i.negatives,
			i.overall_assessment,
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
		AND i.id > $3
		GROUP BY 
			i.id,
			i.interview_type,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.description,
			i.interviewers_decision,
			i.positives,
			i.negatives,
			i.overall_assessment,
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

	rows, err := p.pool.Query(ctx, query,
		req.OpeningID,      // $1
		orgUser.EmployerID, // $2
		req.PaginationKey,  // $3
		req.Limit,          // $4
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no interviews found", "opening_id", req.OpeningID)
			return []employer.EmployerInterview{}, nil
		}

		p.log.Err("failed to get interviews by opening", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	interviews := []employer.EmployerInterview{}
	for rows.Next() {
		var interview employer.EmployerInterview
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

func (p *PG) GetEmployerInterviewsByCandidacy(
	ctx context.Context,
	req employer.GetEmployerInterviewsByCandidacyRequest,
) ([]employer.EmployerInterview, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
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
			i.positives,
			i.negatives,
			i.overall_assessment,
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
		LEFT JOIN interview_interviewers ii ON i.id = ii.interview_id
		LEFT JOIN org_users ou ON ii.interviewer_id = ou.id
		LEFT JOIN org_users fb ON i.feedback_submitted_by = fb.id
		WHERE i.candidacy_id = $1
		AND i.employer_id = $2
		GROUP BY 
			i.id,
			i.interview_type,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.description,
			i.interviewers_decision,
			i.positives,
			i.negatives,
			i.overall_assessment,
			i.feedback_to_candidate,
			i.created_at,
			fb.id,
			fb.name,
			fb.email,
			i.feedback_submitted_at
		ORDER BY i.start_time ASC, i.id ASC
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

	rows, err := p.pool.Query(ctx, query,
		req.CandidacyID,    // $1
		orgUser.EmployerID, // $2
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no interviews found", "candidacy_id", req.CandidacyID)
			return []employer.EmployerInterview{}, nil
		}

		p.log.Err("failed to get interviews by candidacy", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	interviews := []employer.EmployerInterview{}
	for rows.Next() {
		var interview employer.EmployerInterview
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

	p.log.Dbg("found interviews", "interviews", interviews)
	return interviews, nil
}
