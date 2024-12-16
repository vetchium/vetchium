package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
)

func (p *PG) AddInterviewer(
	ctx context.Context,
	addInterviewerReq db.AddInterviewerRequest,
) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(ctx)

	// Insert into interview_interviewers table

	// Insert email
	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
    VALUES ($1, $2, $3, $4, $5, $6),
    ($7, $8, $9, $10, $11, $12)
`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		addInterviewerReq.InterviewerNotificationEmail.EmailFrom,
		addInterviewerReq.InterviewerNotificationEmail.EmailTo,
		addInterviewerReq.InterviewerNotificationEmail.EmailSubject,
		addInterviewerReq.InterviewerNotificationEmail.EmailHTMLBody,
		addInterviewerReq.InterviewerNotificationEmail.EmailTextBody,
		addInterviewerReq.InterviewerNotificationEmail.EmailState,
		// -- End of interviewer notification email
		addInterviewerReq.WatcherNotificationEmail.EmailFrom,
		addInterviewerReq.WatcherNotificationEmail.EmailTo,
		addInterviewerReq.WatcherNotificationEmail.EmailSubject,
		addInterviewerReq.WatcherNotificationEmail.EmailHTMLBody,
		addInterviewerReq.WatcherNotificationEmail.EmailTextBody,
		addInterviewerReq.WatcherNotificationEmail.EmailState,
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

func (p *PG) GetWatchersInfoByInterviewID(
	ctx context.Context,
	interviewID string,
) (db.WatchersInfo, error) {
	return db.WatchersInfo{}, nil
}
