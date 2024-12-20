package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

const (
	validationOK                  = "OK"
	validationInterviewNotFound   = "INTERVIEW_NOT_FOUND"
	validationInvalidInterview    = "INVALID_INTERVIEW_STATE"
	validationInterviewerNotFound = "INTERVIEWER_NOT_FOUND"
	validationInvalidInterviewer  = "INVALID_INTERVIEWER_STATE"
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

	interviewersQuery := `
WITH interview_check AS (
	SELECT i.id, i.interview_state, i.employer_id
	FROM interviews i
	WHERE i.id = $1
),
interviewer_check AS (
	SELECT ou.id, ou.org_user_state 
	FROM org_users ou, interview_check ic
	WHERE ou.email = $2 
	AND ou.employer_id = ic.employer_id
),
state_validation AS (
	SELECT 
		CASE 
			WHEN ic.id IS NULL THEN $5
			WHEN ic.interview_state != $3 THEN $6
			WHEN iw.id IS NULL THEN $7
			WHEN NOT (iw.org_user_state = ANY($4::org_user_states[])) THEN $8
			ELSE $9
		END as validation_result
	FROM interview_check ic
	LEFT JOIN interviewer_check iw ON true
)
INSERT INTO interview_interviewers (interview_id, interviewer_id, employer_id, rsvp_status)
SELECT 
	$1,
	iw.id,
	ic.employer_id,
	$10::rsvp_status
FROM interview_check ic, interviewer_check iw
WHERE EXISTS (
	SELECT 1 FROM state_validation WHERE validation_result = $9
)
RETURNING (SELECT validation_result FROM state_validation)`

	var validationResult string
	err = tx.QueryRow(
		ctx,
		interviewersQuery,
		addInterviewerReq.InterviewID,
		addInterviewerReq.InterviewerEmailAddr,
		string(common.ScheduledInterviewState),
		[]string{
			string(employer.ActiveOrgUserState),
			string(employer.AddedOrgUserState),
			string(employer.ReplicatedOrgUserState),
		},
		validationInterviewNotFound,
		validationInvalidInterview,
		validationInterviewerNotFound,
		validationInvalidInterviewer,
		validationOK,
		common.NotSetRSVP,
	).Scan(&validationResult)

	if err != nil || validationResult != validationOK {
		switch validationResult {
		case validationInterviewNotFound:
			return db.ErrNoInterview
		case validationInvalidInterview:
			return db.ErrInvalidInterviewState
		case validationInterviewerNotFound:
			return db.ErrNoOrgUser
		case validationInvalidInterviewer:
			return db.ErrInterviewerNotActive
		default:
			p.log.Err("failed to insert interviewer", "error", err)
			return db.ErrInternal
		}
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
		WHERE i.id = $2
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
