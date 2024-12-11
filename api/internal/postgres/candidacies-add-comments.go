package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) AddEmployerCandidacyComment(
	ctx context.Context,
	empCommentReq employer.AddEmployerCandidacyCommentRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	// Check if the user has the 'Admin' role
	isAdmin := false
	for _, role := range orgUser.OrgUserRoles {
		if role == common.Admin {
			isAdmin = true
			break
		}
	}

	query := `
	WITH valid_user AS (
		SELECT 1 AS is_authorized
		FROM openings o
		LEFT JOIN opening_watchers ow ON o.employer_id = ow.employer_id AND o.id = ow.opening_id
		WHERE o.employer_id = $1
		  AND o.id = (SELECT opening_id FROM candidacies WHERE id = $2)
		  AND (
			  $3 = o.hiring_manager OR 
			  $3 = o.recruiter OR 
			  $3 = ow.watcher_id OR 
			  $4 = true
		  )
	),
	valid_candidacy AS (
		SELECT 1 AS is_valid_state
		FROM candidacies c
		WHERE c.id = $2
		  AND c.candidacy_state = ANY($5)
	)
	INSERT INTO candidacy_comments (id, author_type, org_user_id, comment_text, candidacy_id, employer_id, created_at)
	SELECT gen_random_uuid(), 'ORG_USER', $3, $6, $2, $1, timezone('UTC', now())
	WHERE EXISTS (SELECT 1 FROM valid_user)
	  AND EXISTS (SELECT 1 FROM valid_candidacy)
	RETURNING id, (SELECT is_authorized FROM valid_user), (SELECT is_valid_state FROM valid_candidacy);
	`

	var (
		commentID    uuid.UUID
		isAuthorized sql.NullBool
		isValidState sql.NullBool
	)
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUser.EmployerID,
		empCommentReq.CandidacyID,
		orgUser.ID,
		isAdmin,
		[]string{
			string(common.InterviewingCandidacyState),
			string(common.OfferedCandidacyState),
		},
		empCommentReq.Comment,
	).Scan(&commentID, &isAuthorized, &isValidState)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if !isAuthorized.Bool {
				return uuid.UUID{}, db.ErrUnauthorizedComment
			}
			if !isValidState.Bool {
				return uuid.UUID{}, db.ErrInvalidCandidacyState
			}
			return uuid.UUID{}, db.ErrNoOpening
		}
		p.log.Err("failed to add employer candidacy comment", "error", err)
		return uuid.UUID{}, db.ErrInternal
	}

	return commentID, nil
}

func (p *PG) AddHubCandidacyComment(
	ctx context.Context,
	hubCommentReq hub.AddHubCandidacyCommentRequest,
) (uuid.UUID, error) {
	const (
		errorCandidacyNotFound = "candidacy_not_found"
		errorInvalidState      = "invalid_state"
		errorValidState        = "valid_state"
		errorUnauthorized      = "unauthorized"
	)

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
WITH candidacy_check AS (
    SELECT
        c.candidacy_state,
        c.hub_user_id,
        CASE WHEN c.id IS NULL THEN
            $3
        WHEN c.candidacy_state NOT IN ($8, $9) THEN
            $4
        ELSE
            $5
        END AS error_code
    FROM
        candidacies c
    WHERE
        c.id = $1
),
auth_check AS (
    SELECT
        cc.*,
        cc.hub_user_id = $2 AS is_authorized
    FROM
        candidacy_check cc
),
insert_comment AS (
    INSERT INTO candidacy_comments (candidacy_id, hub_user_id, author_type, comment_text)
    SELECT
        $1,
        ac.hub_user_id,
        $6,
        $7
    FROM
        auth_check ac
    WHERE
        ac.is_authorized
        AND ac.error_code = $5
    RETURNING
        id,
        ac.error_code,
        ac.is_authorized
)
SELECT
    id,
    error_code,
    is_authorized
FROM
    insert_comment
UNION ALL
SELECT
    NULL::uuid,
    ac.error_code,
    ac.is_authorized
FROM
    auth_check ac
WHERE
    NOT EXISTS (
        SELECT
            1
        FROM
            insert_comment)
LIMIT 1
`

	var (
		commentID    uuid.UUID
		errorCode    string
		isAuthorized bool
	)
	err := p.pool.QueryRow(
		ctx,
		query,
		hubCommentReq.CandidacyID,
		hubUser.ID,
		errorCandidacyNotFound,
		errorInvalidState,
		errorValidState,
		db.HubUserAuthorType,
		hubCommentReq.Comment,
		common.InterviewingCandidacyState,
		common.OfferedCandidacyState,
	).Scan(&commentID, &errorCode, &isAuthorized)

	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.UUID{}, db.ErrNoOpening
		}
		p.log.Err("failed to add hub candidacy comment", "error", err)
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
