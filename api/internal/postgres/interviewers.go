package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) AddInterviewers(
	ctx context.Context,
	addInterviewersReq db.AddInterviewersRequest,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	// Combined verification and insertion query
	insertQuery := `
WITH active_interviewers AS (
	SELECT id, email
	FROM org_users
	WHERE email = ANY($2::text[])
	AND org_user_state IN ('ACTIVE_ORG_USER', 'REPLICATED_ORG_USER')
),
verification AS (
	SELECT COUNT(*) = array_length($2::text[], 1) as all_active
	FROM active_interviewers
),
insertion AS (
	INSERT INTO interview_interviewers (interview_id, interviewer_id, employer_id, rsvp_status)
	SELECT $1, ai.id, i.employer_id, $3::rsvp_status
	FROM interviews i
	CROSS JOIN active_interviewers ai
	WHERE i.id = $1
	AND (SELECT all_active FROM verification)
	RETURNING 1
)
SELECT EXISTS (SELECT 1 FROM insertion) as inserted
`

	var inserted bool
	err = tx.QueryRow(
		ctx,
		insertQuery,
		addInterviewersReq.InterviewID,
		addInterviewersReq.Interviewers,
		common.NotSetRSVP,
	).Scan(&inserted)
	if err != nil {
		p.log.Err("failed to insert interviewers", "error", err)
		return db.ErrInternal
	}

	if !inserted {
		p.log.Dbg("one or more interviewers not in active state")
		return db.ErrInterviewNotActive
	}

	// Insert email
	emailQuery := `
		INSERT INTO emails (
			email_from,
			email_to,
			email_subject,
			email_html_body,
			email_text_body,
			email_state
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		addInterviewersReq.Email.EmailFrom,
		addInterviewersReq.Email.EmailTo,
		addInterviewersReq.Email.EmailSubject,
		addInterviewersReq.Email.EmailHTMLBody,
		addInterviewersReq.Email.EmailTextBody,
		addInterviewersReq.Email.EmailState,
	)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

func (p *PG) RemoveInterviewer(
	ctx context.Context,
	removeInterviewerReq employer.RemoveInterviewerRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	const (
		statusNoInterview = "no_interview"
		statusWrongState  = "wrong_state"
		statusOK          = "ok"
	)

	query := `
WITH interview_check AS (
    SELECT 
        CASE
            WHEN NOT EXISTS (
                SELECT 1 FROM interviews 
                WHERE id = $1 AND employer_id = $2
            ) THEN $4
            WHEN EXISTS (
                SELECT 1 FROM interviews
                WHERE id = $1 AND employer_id = $2
                AND interview_state != $7
            ) THEN $5
            ELSE $6
        END as status
),
delete_result AS (
    DELETE FROM interview_interviewers
    WHERE interview_id = $1
    AND employer_id = $2
    AND interviewer_id = (
        SELECT id FROM org_users 
        WHERE email = $3 AND employer_id = $2
    )
    AND EXISTS (
        SELECT 1 FROM interview_check 
        WHERE status = $6
    )
)
SELECT status FROM interview_check;
`

	var status string
	err := p.pool.QueryRow(
		ctx,
		query,
		removeInterviewerReq.InterviewID,
		orgUser.EmployerID,
		removeInterviewerReq.OrgUserEmail,
		statusNoInterview,
		statusWrongState,
		statusOK,
		common.ScheduledInterviewState,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to remove interviewer", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNoInterview:
		return db.ErrNoInterview
	case statusWrongState:
		return db.ErrInvalidInterviewState
	case statusOK:
		return nil
	default:
		p.log.Err("unexpected status", "status", status)
		return db.ErrInternal
	}
}
