package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
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
