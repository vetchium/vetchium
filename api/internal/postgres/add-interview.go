package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) AddInterview(
	ctx context.Context,
	req db.AddInterviewRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	const (
		statusNoCandidacy  = "no_candidacy"
		statusInvalidState = "invalid_state"
		statusOK           = "ok"
	)

	query := `
WITH candidacy_check AS (
	SELECT
		CASE
			WHEN NOT EXISTS (
				SELECT 1 FROM candidacies
				WHERE id = $1 AND employer_id = $2
			) THEN $9
			WHEN EXISTS (
				SELECT 1 FROM candidacies
				WHERE id = $1 
				AND employer_id = $2
				AND candidacy_state != $10::candidacy_states
			) THEN $11
			ELSE $12
		END as status
),
interview_insert AS (
	INSERT INTO interviews (
		id,
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
		$14,
		$1,                             -- candidacy_id
		$3::interview_types,            -- interview_type
		$4::interview_states,           -- interview_state
		$5::timestamptz,               -- start_time
		$6::timestamptz,               -- end_time
		$7,                            -- description
		$8,                            -- created_by
		$2,                            -- employer_id
		$13::rsvp_status               -- candidate_rsvp default
	WHERE (SELECT status FROM candidacy_check) = $12
	RETURNING id
)
SELECT
	COALESCE(
		(SELECT id::text FROM interview_insert),
		(SELECT status FROM candidacy_check)
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
		statusInvalidState,                // $11
		statusOK,                          // $12
		common.NotSetRSVP,                 // $13
		req.InterviewID,                   // $14
	).Scan(&result)
	if err != nil {
		p.log.Err("failed to add interview", "error", err)
		return db.ErrInternal
	}

	switch result {
	case statusNoCandidacy:
		return db.ErrNoCandidacy
	case statusInvalidState:
		return db.ErrInvalidCandidacyState
	default:
		return nil
	}
}
