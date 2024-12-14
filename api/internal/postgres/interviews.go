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
