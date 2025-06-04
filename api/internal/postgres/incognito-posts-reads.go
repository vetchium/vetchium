package postgres

import (
	"context"
	"database/sql"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) GetIncognitoPost(
	ctx context.Context,
	req hub.GetIncognitoPostRequest,
) (hub.IncognitoPost, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return hub.IncognitoPost{}, err
	}

	var post hub.IncognitoPost
	query := `
		SELECT 
			ip.id,
			ip.content,
			ip.created_at
		FROM incognito_posts ip
		WHERE ip.id = $1 AND ip.author_id = $2
	`

	err = pg.pool.QueryRow(ctx, query, req.IncognitoPostID, hubUserID).Scan(
		&post.IncognitoPostID,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return hub.IncognitoPost{}, db.ErrNoIncognitoPost
		}
		return hub.IncognitoPost{}, err
	}

	tags, err := pg.getIncognitoPostTags(ctx, req.IncognitoPostID)
	if err != nil {
		return hub.IncognitoPost{}, err
	}
	post.Tags = tags

	return post, nil
}

func (pg *PG) GetIncognitoPostComments(
	ctx context.Context,
	req hub.GetIncognitoPostCommentsRequest,
) (hub.GetIncognitoPostCommentsResponse, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	if !pg.canAccessIncognitoPost(ctx, req.IncognitoPostID, hubUserID) {
		return hub.GetIncognitoPostCommentsResponse{}, db.ErrNoIncognitoPost
	}

	var comments []hub.IncognitoPostComment
	query := `
		SELECT 
			ipc.id,
			ipc.content,
			ipc.parent_comment_id,
			ipc.depth,
			ipc.created_at,
			ipc.upvotes_count,
			ipc.downvotes_count,
			CASE WHEN ipc.author_id = $2 THEN TRUE ELSE FALSE END as is_created_by_me,
			CASE 
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv 
					WHERE ipcv.comment_id = ipc.id 
					AND ipcv.user_id = $2 
					AND ipcv.vote_value = 1
				) THEN 'UPVOTE'
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv 
					WHERE ipcv.comment_id = ipc.id 
					AND ipcv.user_id = $2 
					AND ipcv.vote_value = -1
				) THEN 'DOWNVOTE'
				ELSE 'NO_VOTE'
			END as my_vote
		FROM incognito_post_comments ipc
		WHERE ipc.incognito_post_id = $1
		ORDER BY 
			CASE WHEN ipc.parent_comment_id IS NULL THEN ipc.created_at END ASC,
			ipc.parent_comment_id ASC NULLS FIRST,
			ipc.created_at ASC
	`

	rows, err := pg.pool.Query(ctx, query, req.IncognitoPostID, hubUserID)
	if err != nil {
		return hub.GetIncognitoPostCommentsResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment hub.IncognitoPostComment
		var parentCommentID sql.NullString

		err := rows.Scan(
			&comment.CommentID,
			&comment.Content,
			&parentCommentID,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.Upvotes,
			&comment.Downvotes,
			&comment.IsCreatedByMe,
			&comment.MyVote,
		)
		if err != nil {
			return hub.GetIncognitoPostCommentsResponse{}, err
		}

		if parentCommentID.Valid {
			comment.InReplyTo = &parentCommentID.String
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	return hub.GetIncognitoPostCommentsResponse{
		Comments: comments,
	}, nil
}

func (pg *PG) getIncognitoPostTags(
	ctx context.Context,
	incognitoPostID string,
) ([]common.VTag, error) {
	query := `
		SELECT t.id, t.display_name
		FROM incognito_post_tags ipt
		JOIN tags t ON ipt.tag_id = t.id
		WHERE ipt.incognito_post_id = $1
		ORDER BY t.display_name
	`

	rows, err := pg.pool.Query(ctx, query, incognitoPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []common.VTag
	for rows.Next() {
		var tag common.VTag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func (pg *PG) canAccessIncognitoPost(
	ctx context.Context,
	incognitoPostID string,
	hubUserID string,
) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM incognito_posts
		WHERE id = $1 AND author_id = $2
	`

	err := pg.pool.QueryRow(ctx, query, incognitoPostID, hubUserID).Scan(&count)
	return err == nil && count > 0
}
