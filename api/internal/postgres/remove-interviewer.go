package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) RemoveInterviewer(
	ctx context.Context,
	removeInterviewerReq db.RemoveInterviewerRequest,
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

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)
	p.log.Dbg("Transaction started")

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
	err = tx.QueryRow(
		ctx,
		query,
		removeInterviewerReq.InterviewID,
		orgUser.EmployerID,
		removeInterviewerReq.RemovedInterviewerEmailAddr,
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
		// Go ahead and send email notification
		p.log.Dbg("interview_interviewers deleted")
	default:
		p.log.Err("unexpected status", "status", status)
		return db.ErrInternal
	}

	candidacyCommentQuery := `
INSERT INTO candidacy_comments (
	author_type,
	org_user_id,
	comment_text,
	candidacy_id,
	employer_id
) VALUES (
	$1, $2, $3, 
	(SELECT candidacy_id FROM interviews WHERE id = $4 AND employer_id = $5),
	$5
)
RETURNING id
`

	var candidacyCommentID string
	err = tx.QueryRow(
		ctx,
		candidacyCommentQuery,
		db.OrgUserAuthorType,
		orgUser.ID,
		removeInterviewerReq.CandidacyComment,
		removeInterviewerReq.InterviewID,
		orgUser.EmployerID,
	).Scan(&candidacyCommentID)
	if err != nil {
		p.log.Err("failed to insert candidacy comment", "error", err)
		return db.ErrInternal
	}

	emailQuery := `
INSERT INTO emails (
	email_from,
	email_to,
	email_subject,
	email_html_body,
	email_text_body,
	email_state
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING email_key
`
	var emailKey string
	err = tx.QueryRow(
		ctx,
		emailQuery,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailFrom,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailTo,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailSubject,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailHTMLBody,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailTextBody,
		removeInterviewerReq.RemovedInterviewerEmailNotification.EmailState,
	).Scan(&emailKey)
	if err != nil {
		p.log.Err("failed to insert email", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}
	p.log.Dbg("Transaction committed")

	return nil
}
