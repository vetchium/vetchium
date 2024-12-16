package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
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

	var state struct {
		candidacyState  string
		orgUserState    string
		interviewExists bool
	}

	// Insert into interview_interviewers table with state checks
	err = tx.QueryRow(ctx, `
		WITH valid_states AS (
			SELECT i.interview_id, 
				   c.candidacy_state, 
				   ou.org_user_state,
				   true as interview_exists
			FROM interviews i
			LEFT JOIN candidacies c ON i.candidacy_id = c.candidacy_id
			LEFT JOIN org_users ou ON ou.email = $2
			WHERE i.interview_id = $1
		)
		INSERT INTO interview_interviewers (
			interview_id,
			interviewer_email,
			created_at
		) 
		SELECT 
			$1, 
			$2, 
			NOW()
		FROM valid_states
		WHERE candidacy_state = $3
		AND org_user_state IN ($4, $5, $6)
		RETURNING 
			(SELECT candidacy_state FROM valid_states),
			(SELECT org_user_state FROM valid_states),
			(SELECT interview_exists FROM valid_states)
	`,
		addInterviewerReq.InterviewID,
		addInterviewerReq.InterviewerEmailAddr,
		common.InterviewingCandidacyState,
		employer.ActiveOrgUserState,
		employer.AddedOrgUserState,
		employer.ReplicatedOrgUserState,
	).Scan(&state.candidacyState, &state.orgUserState, &state.interviewExists)

	if err != nil {
		p.log.Err("failed to insert interviewer", "error", err)
		if !state.interviewExists {
			return db.ErrNoInterview
		}
		if state.orgUserState == "" {
			return db.ErrNoOrgUser
		}
		if state.orgUserState != string(employer.ActiveOrgUserState) &&
			state.orgUserState != string(employer.AddedOrgUserState) &&
			state.orgUserState != string(employer.ReplicatedOrgUserState) {
			return db.ErrNoOrgUser
		}
		return db.ErrInvalidCandidacyState
	}

	// Insert into candidacy_comments table
	_, err = tx.Exec(ctx, `
		INSERT INTO candidacy_comments (
			candidacy_id,
			comment_text,
			created_at
		) SELECT 
			c.candidacy_id,
			$1,
			NOW()
		FROM interviews i
		JOIN candidacies c ON i.candidacy_id = c.candidacy_id
		WHERE i.interview_id = $2
	`, addInterviewerReq.CandidacyComment, addInterviewerReq.InterviewID)
	if err != nil {
		p.log.Err("failed to insert comment", "error", err)
		return db.ErrInternal
	}

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
