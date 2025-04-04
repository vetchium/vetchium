package postgres

import (
	"context"
	"encoding/json"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
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
			i.candidate_rsvp,
			COALESCE(
				jsonb_agg(
					jsonb_build_object(
						'name', ou.name,
						'rsvp_status', ii.rsvp_status
					) 
					ORDER BY ou.name
				) FILTER (WHERE ou.name IS NOT NULL),
				'[]'::jsonb
			) as interviewer_data
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
			i.description,
			i.candidate_rsvp
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
		var interviewerData []byte
		if err := rows.Scan(
			&interview.InterviewID,
			&interview.InterviewState,
			&interview.StartTime,
			&interview.EndTime,
			&interview.InterviewType,
			&interview.Description,
			&interview.CandidateRSVP,
			&interviewerData,
		); err != nil {
			p.log.Err("failed to scan hub interview", "error", err)
			return nil, db.ErrInternal
		}

		// Parse the interviewer data from JSON
		var interviewers []struct {
			Name       string            `json:"name"`
			RSVPStatus common.RSVPStatus `json:"rsvp_status"`
		}
		if err := json.Unmarshal(interviewerData, &interviewers); err != nil {
			p.log.Err("failed to unmarshal interviewer data", "error", err)
			return nil, db.ErrInternal
		}

		// Convert to HubInterviewer slice
		interview.Interviewers = make([]hub.HubInterviewer, len(interviewers))
		for i, interviewer := range interviewers {
			interview.Interviewers[i] = hub.HubInterviewer{
				Name:       interviewer.Name,
				RSVPStatus: interviewer.RSVPStatus,
			}
		}

		interviews = append(interviews, interview)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	return interviews, nil
}
