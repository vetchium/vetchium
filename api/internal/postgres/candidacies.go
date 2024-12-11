package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) AddEmployerCandidacyComment(
	ctx context.Context,
	req employer.AddEmployerCandidacyCommentRequest,
) (uuid.UUID, error) {
	const (
		errorCandidacyNotFound = "candidacy_not_found"
		errorInvalidState      = "invalid_state"
		errorValidState        = "valid_state"
		statusUnauthorized     = "unauthorized"
	)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
		WITH candidacy_check AS (
			SELECT 
				c.candidacy_state, 
				c.employer_id, 
				c.opening_id,
				CASE 
					WHEN c.id IS NULL THEN $5
					WHEN c.candidacy_state NOT IN ('INTERVIEWING', 'OFFERED') THEN $6
					ELSE $7
				END as error_code
			FROM candidacies c
			WHERE c.id = $1
		),
		auth_check AS (
			SELECT 
				cc.*,
				EXISTS (
					SELECT 1
					FROM openings o 
					WHERE o.employer_id = cc.employer_id 
					AND o.id = cc.opening_id
					AND (
						o.recruiter = $3
						OR o.hiring_manager = $3
						OR EXISTS (
							SELECT 1 FROM opening_watchers w
							WHERE w.employer_id = cc.employer_id
							AND w.opening_id = cc.opening_id
							AND w.watcher_id = $3
						)
					)
				) as is_authorized
			FROM candidacy_check cc
			WHERE cc.employer_id = $2
		),
		insert_comment AS (
			INSERT INTO candidacy_comments (
				candidacy_id,
				employer_id,
				author_type,
				org_user_id,
				comment_text
			)
			SELECT 
				$1,
				ac.employer_id,
				$4,
				$3,
				$8
			FROM auth_check ac
			WHERE ac.is_authorized AND ac.error_code = $7
			RETURNING id, ac.error_code, ac.is_authorized
		)
		SELECT id, error_code, is_authorized 
		FROM insert_comment
		UNION ALL
		SELECT NULL::uuid, ac.error_code, ac.is_authorized
		FROM auth_check ac
		WHERE NOT EXISTS (SELECT 1 FROM insert_comment)
		LIMIT 1;
	`

	var (
		commentID    uuid.UUID
		errorCode    string
		isAuthorized bool
	)
	err := p.pool.QueryRow(
		ctx,
		query,
		req.CandidacyID,
		orgUser.EmployerID,
		orgUser.ID,
		db.OrgUserAuthorType,
		errorCandidacyNotFound,
		errorInvalidState,
		errorValidState,
		req.Comment,
	).Scan(&commentID, &errorCode, &isAuthorized)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, db.ErrNoOpening
		}
		p.log.Err("failed to add candidacy comment", "error", err)
		return uuid.UUID{}, db.ErrInternal
	}

	// If comment wasn't inserted, determine why based on error_code and is_authorized
	if commentID == uuid.Nil {
		switch errorCode {
		case errorCandidacyNotFound:
			return uuid.UUID{}, db.ErrNoOpening
		case errorInvalidState:
			return uuid.UUID{}, db.ErrInvalidCandidacyState
		default:
			if !isAuthorized {
				return uuid.UUID{}, db.ErrUnauthorizedComment
			}
			return uuid.UUID{}, db.ErrInternal
		}
	}

	return commentID, nil
}
