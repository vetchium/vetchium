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
    SELECT
        1 AS is_authorized
    FROM
        openings o
        LEFT JOIN opening_watchers ow ON o.employer_id = ow.employer_id
            AND o.id = ow.opening_id
    WHERE
        o.employer_id = $1
        AND o.id = (
            SELECT
                opening_id
            FROM
                candidacies
            WHERE
                id = $2)
            AND ($3 = o.hiring_manager
                OR $3 = o.recruiter
                OR $3 = ow.watcher_id
                OR $4 = TRUE)
),
valid_candidacy AS (
    SELECT
        1 AS is_valid_state
    FROM
        candidacies c
    WHERE
        c.id = $2
        AND c.candidacy_state = ANY ($5))
INSERT INTO candidacy_comments (author_type, org_user_id, comment_text, candidacy_id, employer_id, created_at)
SELECT
    $6,
    $3,
    $7,
    $2,
    $1,
    timezone('UTC', now())
WHERE
    EXISTS (
        SELECT
            1
        FROM
            valid_user)
    AND EXISTS (
        SELECT
            1
        FROM
            valid_candidacy)
RETURNING
    id,
    (
        SELECT
            is_authorized
        FROM
            valid_user),
    (
        SELECT
            is_valid_state
        FROM
            valid_candidacy)
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
		db.OrgUserAuthorType,
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
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
WITH valid_candidacy AS (
    SELECT
        1 AS is_valid_state
    FROM
        candidacies c
        JOIN applications a ON c.application_id = a.id
    WHERE
        c.id = $1
        AND c.candidacy_state = ANY ($4)
        AND a.hub_user_id = $2
),
valid_candidacy_id AS (
    SELECT
        1 AS is_valid_id
    FROM
        candidacies
    WHERE
        id = $1
)
INSERT INTO candidacy_comments (author_type, hub_user_id, comment_text, candidacy_id, employer_id, created_at)
SELECT
    $3,
    $2,
    $5,
    $1,
    c.employer_id,
    timezone('UTC', now())
FROM
    candidacies c
WHERE
    c.id = $1
    AND EXISTS (
        SELECT
            1
        FROM
            valid_candidacy)
    AND EXISTS (
        SELECT
            1
        FROM
            valid_candidacy_id)
RETURNING
    id,
    (
        SELECT
            is_valid_state
        FROM
            valid_candidacy),
    (
        SELECT
            is_valid_id
        FROM
            valid_candidacy_id)
`

	var (
		commentID    uuid.UUID
		isValidState sql.NullBool
		isValidID    sql.NullBool
	)
	err := p.pool.QueryRow(
		ctx,
		query,
		hubCommentReq.CandidacyID,
		hubUser.ID,
		db.HubUserAuthorType,
		[]string{
			string(common.InterviewingCandidacyState),
			string(common.OfferedCandidacyState),
		},
		hubCommentReq.Comment,
	).Scan(&commentID, &isValidState, &isValidID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if !isValidID.Bool {
				return uuid.UUID{}, db.ErrUnauthorizedComment
			}
			if !isValidState.Bool {
				return uuid.UUID{}, db.ErrInvalidCandidacyState
			}
			return uuid.UUID{}, db.ErrNoOpening
		}
		p.log.Err("failed to add hub candidacy comment", "error", err)
		return uuid.UUID{}, db.ErrInternal
	}

	return commentID, nil
}
