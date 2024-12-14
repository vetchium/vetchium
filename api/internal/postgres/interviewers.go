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

	// First verify all interviewers exist and are in active state
	verifyQuery := `
		SELECT COUNT(*) = 0
		FROM (
			SELECT id 
			FROM org_users 
			WHERE id = ANY($1::uuid[])
			AND org_user_state NOT IN ('ACTIVE_ORG_USER', 'REPLICATED_ORG_USER')
		) inactive_users
	`

	var allActive bool
	err = tx.QueryRow(ctx, verifyQuery, addInterviewersReq.OrgUserIDs).
		Scan(&allActive)
	if err != nil {
		p.log.Err("failed to verify interviewer states", "error", err)
		return db.ErrInternal
	}

	if !allActive {
		p.log.Dbg("one or more interviewers not in active state")
		return db.ErrInterviewNotActive
	}

	// Insert interviewers
	insertQuery := `
		INSERT INTO interview_interviewers 
		(interview_id, interviewer_id, employer_id, rsvp_status)
		SELECT $1, unnest($2::uuid[]), i.employer_id, $3::rsvp_status
		FROM interviews i
		WHERE i.id = $1
	`

	_, err = tx.Exec(
		ctx,
		insertQuery,
		addInterviewersReq.InterviewID,
		addInterviewersReq.OrgUserIDs,
		common.NotSetRSVP,
	)
	if err != nil {
		p.log.Err("failed to insert interviewers", "error", err)
		return db.ErrInternal
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
