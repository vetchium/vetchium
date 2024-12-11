package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
)

func (p *PG) GetEmployerCandidacyComments(
	ctx context.Context,
	empGetCommentsReq common.GetCandidacyCommentsRequest,
) ([]common.CandidacyComment, error) {
	var candidacyComments []common.CandidacyComment

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return candidacyComments, db.ErrInternal
	}

	query := `
WITH access_check AS (
	SELECT EXISTS (
		SELECT 1 FROM candidacies 
		WHERE id = $1 AND employer_id = $2
	)
)
SELECT 
	cc.comment_id,
	COALESCE(ou.name, hu.full_name) as commenter_name,
	cc.author_type,
	cc.comment_text as content,
	cc.created_at
FROM candidacy_comments cc
LEFT JOIN org_users ou ON cc.org_user_id = ou.id
LEFT JOIN hub_users hu ON cc.hub_user_id = hu.id,
access_check
WHERE cc.candidacy_id = $1
AND access_check
ORDER BY cc.created_at DESC
`

	rows, err := p.pool.Query(ctx, query,
		empGetCommentsReq.CandidacyID,
		orgUser.EmployerID,
		db.OrgUserAuthorType,
		db.HubUserAuthorType,
	)
	if err != nil {
		p.log.Err("failed to query candidacy comments", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var comment common.CandidacyComment
		err := rows.Scan(
			&comment.CommentID,
			&comment.CommenterName,
			&comment.CommenterType,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy comment", "error", err)
			return nil, db.ErrInternal
		}

		candidacyComments = append(candidacyComments, comment)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	return candidacyComments, nil
}

func (p *PG) GetHubCandidacyComments(
	ctx context.Context,
	hubGetCommentsReq common.GetCandidacyCommentsRequest,
) ([]common.CandidacyComment, error) {
	var candidacyComments []common.CandidacyComment

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("failed to get hubUser from context")
		return candidacyComments, db.ErrInternal
	}

	query := `
SELECT 
	cc.comment_id,
	COALESCE(ou.name, hu.full_name) as commenter_name,
	cc.author_type,
	cc.comment_text as content,
	cc.created_at
FROM candidacy_comments cc
WHERE cc.candidacy_id = $1
AND cc.hub_user_id = $2
ORDER BY cc.created_at DESC
`

	rows, err := p.pool.Query(ctx, query,
		hubGetCommentsReq.CandidacyID,
		hubUser.ID,
	)
	if err != nil {
		p.log.Err("failed to query candidacy comments", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var comment common.CandidacyComment
		err := rows.Scan(
			&comment.CommentID,
			&comment.CommenterName,
			&comment.CommenterType,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			p.log.Err("failed to scan candidacy comment", "error", err)
			return nil, db.ErrInternal
		}

		candidacyComments = append(candidacyComments, comment)
	}

	if err = rows.Err(); err != nil {
		p.log.Err("error iterating over rows", "error", err)
		return nil, db.ErrInternal
	}

	return candidacyComments, nil
}
