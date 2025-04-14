package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
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
	FROM (SELECT 1) AS dummy
	LEFT JOIN interview_check ic ON true
	LEFT JOIN interviewer_check iw ON true
),
insert_interviewer AS (
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
	RETURNING true
)
SELECT validation_result FROM state_validation`

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
			$3,
			$4,
			$1,
			timezone('UTC', now()),
			c.id,
			c.employer_id
		FROM interviews i
		JOIN candidacies c ON i.candidacy_id = c.id
		WHERE i.id = $2
	`,
		addInterviewerReq.CandidacyComment,
		addInterviewerReq.InterviewID,
		string(db.OrgUserAuthorType),
		orgUser.ID,
	)
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

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

func (p *PG) GetStakeholdersByInterview(
	ctx context.Context,
	interviewID string,
) (db.Stakeholders, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.Stakeholders{}, db.ErrInternal
	}

	query := `
	WITH interview_data AS (
		SELECT
			c.employer_id,
			c.opening_id
		FROM interviews i
		JOIN candidacies c ON i.candidacy_id = c.id
		WHERE i.id = $1
		AND c.employer_id = $2
	)
	SELECT
		o.id as opening_id,
		o.title as opening_title,
		o.state as opening_state,
		o.opening_type,
		hm.name as hiring_manager_name,
		hm.email as hiring_manager_email,
		r.name as recruiter_name,
		r.email as recruiter_email,
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'name', w.name,
					'email', w.email
				)
				ORDER BY w.name
			) FILTER (WHERE w.id IS NOT NULL),
			'[]'::jsonb
		) as watchers
	FROM interview_data id
	JOIN openings o ON id.opening_id = o.id AND id.employer_id = o.employer_id
	JOIN org_users hm ON o.hiring_manager = hm.id
	JOIN org_users r ON o.recruiter = r.id
	LEFT JOIN opening_watchers ow ON o.id = ow.opening_id AND o.employer_id = ow.employer_id
	LEFT JOIN org_users w ON ow.watcher_id = w.id
	GROUP BY
		o.id,
		o.title,
		o.state,
		o.opening_type,
		hm.name,
		hm.email,
		r.name,
		r.email
	`

	var stakeholders db.Stakeholders
	var watchersData []byte

	err := p.pool.QueryRow(ctx, query, interviewID, orgUser.EmployerID).Scan(
		&stakeholders.OpeningID,
		&stakeholders.OpeningTitle,
		&stakeholders.OpeningState,
		&stakeholders.OpeningType,
		&stakeholders.HiringManager.Name,
		&stakeholders.HiringManager.Email,
		&stakeholders.Recruiter.Name,
		&stakeholders.Recruiter.Email,
		&watchersData,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("not found or inaccessible", "interview_id", interviewID)
			return db.Stakeholders{}, db.ErrNoInterview
		}

		p.log.Err("stakeholders", "error", err, "interview_id", interviewID)
		return db.Stakeholders{}, db.ErrInternal
	}

	// Parse watchers from JSON
	var watchers []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err = json.Unmarshal(watchersData, &watchers)
	if err != nil {
		p.log.Err("failed to unmarshal watchers data", "error", err)
		return db.Stakeholders{}, db.ErrInternal
	}

	// Convert to OrgUserShort slice
	stakeholders.Watchers = make([]employer.OrgUserShort, len(watchers))
	for i, watcher := range watchers {
		stakeholders.Watchers[i] = employer.OrgUserShort{
			Name:  watcher.Name,
			Email: watcher.Email,
		}
	}

	p.log.Dbg("got stakeholders", "stakeholders", stakeholders)
	return stakeholders, nil
}
