package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) OfferToCandidate(
	ctx context.Context,
	request db.OfferToCandidateReq,
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

	query := `
UPDATE candidacies
SET candidacy_state = $1
WHERE id = $2
AND employer_id = $3
AND candidacy_state = $4
RETURNING id
`

	err = tx.QueryRow(
		ctx,
		query,
		common.OfferedCandidacyState,
		request.CandidacyID,
		orgUser.EmployerID,
		common.InterviewingCandidacyState,
	).Scan(&request.CandidacyID)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Dbg("not found", "candidacy_id", request.CandidacyID)
			return db.ErrNoCandidacy
		}
		p.log.Err("failed to update candidacy status", "error", err)
		return db.ErrInternal
	}

	interviewUpdateQuery := `
UPDATE interviews
SET interview_state = $1
WHERE candidacy_id = $2
AND interview_state = $3
`
	_, err = tx.Exec(
		ctx,
		interviewUpdateQuery,
		common.CancelledInterviewState,
		request.CandidacyID,
		common.ScheduledInterviewState,
	)
	if err != nil {
		p.log.Err("failed to update interview state", "error", err)
		return db.ErrInternal
	}

	commentQuery := `
INSERT INTO candidacy_comments (
	author_type,
	org_user_id,
	comment_text,
	candidacy_id,
	employer_id,
	created_at
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	timezone('UTC', now())
)
`
	_, err = tx.Exec(
		ctx,
		commentQuery,
		db.OrgUserAuthorType,
		orgUser.ID,
		request.Comment,
		request.CandidacyID,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to add comment", "error", err)
		return db.ErrInternal
	}

	emailQuery := `
INSERT INTO emails (
	email_from,
	email_to,
	email_cc,
	email_subject,
	email_html_body,
	email_text_body,
	email_state
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
`
	_, err = tx.Exec(
		ctx,
		emailQuery,
		request.Email.EmailFrom,
		request.Email.EmailTo,
		request.Email.EmailCC,
		request.Email.EmailSubject,
		request.Email.EmailHTMLBody,
		request.Email.EmailTextBody,
		request.Email.EmailState,
	)
	if err != nil {
		p.log.Err("failed to create email", "error", err)
		return db.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

func (p *PG) GetCandidateInfo(
	ctx context.Context,
	candidacyID string,
) (db.CandidateInfo, error) {
	query := `
SELECT
	h.full_name,
	h.email,
	e.company_name,
	o.title
FROM
	candidacies c
	JOIN applications a ON c.application_id = a.id
	JOIN hub_users h ON a.hub_user_id = h.id
	JOIN employers e ON c.employer_id = e.id
	JOIN openings o ON c.employer_id = o.employer_id AND c.opening_id = o.id
WHERE
	c.id = $1
`
	var info db.CandidateInfo
	err := p.pool.QueryRow(ctx, query, candidacyID).Scan(
		&info.CandidateName,
		&info.CandidateEmail,
		&info.CompanyName,
		&info.OpeningTitle,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.CandidateInfo{}, db.ErrNoCandidacy
		}
		p.log.Err("failed to get candidate info", "error", err)
		return db.CandidateInfo{}, db.ErrInternal
	}

	return info, nil
}
