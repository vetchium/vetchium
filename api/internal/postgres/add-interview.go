package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
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

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	const (
		statusNoCandidacy  = "no_candidacy"
		statusInvalidState = "invalid_state"
		statusOK           = "ok"
	)

	// First validate all interviewers exist and are in valid state
	interviewerValidationQuery := `
WITH valid_states AS (
	SELECT ARRAY[
		$1::org_user_states,
		$2::org_user_states,
		$3::org_user_states
	] as states
)
SELECT COUNT(*) = $4::int
FROM (
	SELECT DISTINCT email 
	FROM unnest($5::text[]) as emails(email)
) as unique_emails
WHERE EXISTS (
	SELECT 1 
	FROM org_users ou, valid_states vs
	WHERE ou.email = unique_emails.email 
	AND ou.employer_id = $6
	AND ou.org_user_state = ANY(vs.states)
)`

	var allInterviewersValid bool
	err = tx.QueryRow(
		ctx,
		interviewerValidationQuery,
		string(employer.ActiveOrgUserState),
		string(employer.AddedOrgUserState),
		string(employer.ReplicatedOrgUserState),
		len(req.InterviewerEmails),
		req.InterviewerEmails,
		orgUser.EmployerID,
	).Scan(&allInterviewersValid)

	if err != nil {
		p.log.Err("failed to validate interviewers", "error", err)
		return db.ErrInternal
	}

	if !allInterviewersValid {
		return db.ErrNoOrgUser
	}

	// Now add the interview
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
	err = tx.QueryRow(
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
	}

	// Now add the interviewers
	interviewerInsertQuery := `
INSERT INTO interview_interviewers (interview_id, interviewer_id, employer_id, rsvp_status)
SELECT 
	$1,
	ou.id,
	ou.employer_id,
	$3::rsvp_status
FROM org_users ou
WHERE ou.email = ANY($2::text[])
AND ou.employer_id = $4`

	_, err = tx.Exec(
		ctx,
		interviewerInsertQuery,
		req.InterviewID,
		req.InterviewerEmails,
		common.NotSetRSVP,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to add interviewers", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
