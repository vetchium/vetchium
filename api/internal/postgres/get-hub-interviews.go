package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetHubInterviewsByCandidacy(
	ctx context.Context,
	req hub.GetHubInterviewsByCandidacyRequest,
) ([]hub.HubInterview, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("hub user not found in context")
		return nil, db.ErrInternal
	}

	// TODO: Check if the two sql queries can be combined via CTE

	// First check if the candidacy exists and belongs to the user
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM candidacies c
			JOIN applications a ON c.application_id = a.id
			WHERE c.id = $1 AND a.hub_user_id = $2
		)`

	var exists bool
	err := p.pool.QueryRow(ctx, checkQuery, req.CandidacyID, hubUser.ID).
		Scan(&exists)
	if err != nil {
		p.log.Err("failed to check candidacy existence", "error", err)
		return nil, db.ErrInternal
	}

	if !exists {
		return nil, db.ErrNoApplication
	}

	// Now get interviews if any exist
	query := `
		SELECT 
			i.id,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.interview_type,
			i.description,
			COALESCE(
				array_agg(ou.name) FILTER (WHERE ou.name IS NOT NULL),
				ARRAY[]::TEXT[]
			) as interviewer_names
		FROM interviews i
		LEFT JOIN interview_interviewers ii ON i.id = ii.interview_id
		LEFT JOIN org_users ou ON ii.interviewer_id = ou.id
		WHERE i.candidacy_id = $1
		GROUP BY 
			i.id,
			i.interview_state,
			i.start_time,
			i.end_time,
			i.interview_type,
			i.description
		ORDER BY i.start_time ASC`

	rows, err := p.pool.Query(ctx, query, req.CandidacyID)
	if err != nil {
		p.log.Err("failed to get hub interviews by candidacy", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	interviews := []hub.HubInterview{}
	for rows.Next() {
		var interview hub.HubInterview
		var interviewerNames []string
		if err := rows.Scan(
			&interview.InterviewID,
			&interview.InterviewState,
			&interview.StartTime,
			&interview.EndTime,
			&interview.InterviewType,
			&interview.Description,
			&interviewerNames,
		); err != nil {
			p.log.Err("failed to scan hub interview", "error", err)
			return nil, db.ErrInternal
		}
		interview.Interviewers = interviewerNames
		interviews = append(interviews, interview)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	return interviews, nil
}
