package postgres

import (
	"context"
	"strings"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
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
	defer tx.Rollback(context.Background())

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

	// Get applicant's email
	var applicantEmail string
	err = tx.QueryRow(
		ctx,
		`
		SELECT hu.email
		FROM candidacies c
		JOIN applications a ON c.application_id = a.id
		JOIN hub_users hu ON a.hub_user_id = hu.id
		WHERE c.id = $1 AND c.employer_id = $2
		`,
		req.CandidacyID,
		orgUser.EmployerID,
	).Scan(&applicantEmail)
	if err != nil {
		p.log.Err("failed to get applicant email", "error", err)
		return db.ErrInternal
	}

	// Get interviewer names if any
	var interviewerNames []string
	if len(req.InterviewerEmails) > 0 {
		rows, err := tx.Query(
			ctx,
			`
			SELECT name
			FROM org_users
			WHERE email = ANY($1::text[])
			AND employer_id = $2
			ORDER BY name
			`,
			req.InterviewerEmails,
			orgUser.EmployerID,
		)
		if err != nil {
			p.log.Err("failed to get interviewer names", "error", err)
			return db.ErrInternal
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				p.log.Err("failed to scan interviewer name", "error", err)
				return db.ErrInternal
			}
			interviewerNames = append(interviewerNames, name)
		}
		if err := rows.Err(); err != nil {
			p.log.Err("failed to iterate over interviewer names", "error", err)
			return db.ErrInternal
		}
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

	// Insert candidacy comment
	if req.CandidacyComment != "" {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO candidacy_comments (
				author_type,
				org_user_id,
				comment_text,
				created_at,
				candidacy_id,
				employer_id
			) SELECT 
				$1,
				$2,
				$3,
				timezone('UTC', now()),
				$4,
				$5
			`,
			string(db.OrgUserAuthorType),
			orgUser.ID,
			req.CandidacyComment,
			req.CandidacyID,
			orgUser.EmployerID,
		)
		if err != nil {
			p.log.Err("failed to insert comment", "error", err)
			return db.ErrInternal
		}
	}

	// Insert email notifications
	emailQuery := `
INSERT INTO emails (email_from, email_to, email_subject, email_html_body, email_text_body, email_state)
    VALUES ($1, $2, $3, $4, $5, $6)`

	// Update applicant notification with actual interviewer names
	if len(interviewerNames) > 0 {
		req.ApplicantNotificationEmail.EmailHTMLBody = strings.Replace(
			req.ApplicantNotificationEmail.EmailHTMLBody,
			"[Interviewer names will be added during the transaction]",
			strings.Join(interviewerNames, ", "),
			1,
		)
		req.ApplicantNotificationEmail.EmailTextBody = strings.Replace(
			req.ApplicantNotificationEmail.EmailTextBody,
			"[Interviewer names will be added during the transaction]",
			strings.Join(interviewerNames, ", "),
			1,
		)
	}

	// Insert applicant notification
	_, err = tx.Exec(
		ctx,
		emailQuery,
		req.ApplicantNotificationEmail.EmailFrom,
		[]string{applicantEmail},
		req.ApplicantNotificationEmail.EmailSubject,
		req.ApplicantNotificationEmail.EmailHTMLBody,
		req.ApplicantNotificationEmail.EmailTextBody,
		req.ApplicantNotificationEmail.EmailState,
	)
	if err != nil {
		p.log.Err("failed to insert applicant email", "error", err)
		return db.ErrInternal
	}

	// Insert interviewer and watcher notifications if there are interviewers
	if len(req.InterviewerEmails) > 0 {
		_, err = tx.Exec(
			ctx,
			emailQuery+", ($7, $8, $9, $10, $11, $12)",
			req.InterviewerNotificationEmail.EmailFrom,
			req.InterviewerNotificationEmail.EmailTo,
			req.InterviewerNotificationEmail.EmailSubject,
			req.InterviewerNotificationEmail.EmailHTMLBody,
			req.InterviewerNotificationEmail.EmailTextBody,
			req.InterviewerNotificationEmail.EmailState,
			// -- End of interviewer notification email
			req.WatcherNotificationEmail.EmailFrom,
			req.WatcherNotificationEmail.EmailTo,
			req.WatcherNotificationEmail.EmailSubject,
			req.WatcherNotificationEmail.EmailHTMLBody,
			req.WatcherNotificationEmail.EmailTextBody,
			req.WatcherNotificationEmail.EmailState,
		)
		if err != nil {
			p.log.Err("failed to insert email", "error", err)
			return db.ErrInternal
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}
