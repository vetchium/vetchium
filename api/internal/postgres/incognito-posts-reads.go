package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) GetIncognitoPost(
	ctx context.Context,
	req hub.GetIncognitoPostRequest,
) (hub.IncognitoPost, error) {
	pg.log.Dbg("entered GetIncognitoPost", "id", req.IncognitoPostID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.IncognitoPost{}, err
	}

	// Single query to get post with tags and author check
	query := `
		SELECT 
			ip.id,
			ip.content,
			ip.created_at,
			CASE WHEN ip.author_id = $2 THEN TRUE ELSE FALSE END as is_created_by_me,
			COALESCE(ip.upvotes_count, 0) as upvotes,
			COALESCE(ip.downvotes_count, 0) as downvotes,
			CASE 
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv 
					WHERE ipv.incognito_post_id = ip.id 
					AND ipv.user_id = $2 
					AND ipv.vote_value = 1
				) THEN 'upvote'
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv 
					WHERE ipv.incognito_post_id = ip.id 
					AND ipv.user_id = $2 
					AND ipv.vote_value = -1
				) THEN 'downvote'
				ELSE NULL
			END as my_vote,
			COALESCE(
				ARRAY_AGG(t.id ORDER BY t.display_name) FILTER (WHERE t.id IS NOT NULL),
				'{}'::text[]
			) as tag_ids,
			COALESCE(
				ARRAY_AGG(t.display_name ORDER BY t.display_name) FILTER (WHERE t.display_name IS NOT NULL),
				'{}'::text[]
			) as tag_names
		FROM incognito_posts ip
		LEFT JOIN incognito_post_tags ipt ON ip.id = ipt.incognito_post_id
		LEFT JOIN tags t ON ipt.tag_id = t.id
		WHERE ip.id = $1 AND ip.is_deleted = FALSE
		GROUP BY ip.id, ip.content, ip.created_at, ip.author_id, ip.upvotes_count, ip.downvotes_count
	`

	var post hub.IncognitoPost
	var tagIDs []string
	var tagNames []string
	var myVote sql.NullString

	err = pg.pool.QueryRow(ctx, query, req.IncognitoPostID, hubUserID).Scan(
		&post.IncognitoPostID,
		&post.Content,
		&post.CreatedAt,
		&post.IsCreatedByMe,
		&post.UpvotesCount,
		&post.DownvotesCount,
		&myVote,
		&tagIDs,
		&tagNames,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			pg.log.Dbg("incognito post not found", "id", req.IncognitoPostID)
			return hub.IncognitoPost{}, db.ErrNoIncognitoPost
		}
		pg.log.Err("failed to query incognito post",
			"error", err,
			"id", req.IncognitoPostID,
		)
		return hub.IncognitoPost{}, err
	}

	// Set the user's vote if exists
	if myVote.Valid {
		post.MeUpvoted = myVote.String == "upvote"
		post.MeDownvoted = myVote.String == "downvote"
	}

	// Build tags array
	post.Tags = make([]common.VTag, len(tagIDs))
	for i := 0; i < len(tagIDs) && i < len(tagNames); i++ {
		post.Tags[i] = common.VTag{
			ID:   common.VTagID(tagIDs[i]),
			Name: common.VTagName(tagNames[i]),
		}
	}

	pg.log.Dbg("fetched incognito post",
		"incognito_post_id", req.IncognitoPostID,
		"author_match", post.IsCreatedByMe)
	return post, nil
}

func (pg *PG) GetIncognitoPostComments(
	ctx context.Context,
	req hub.GetIncognitoPostCommentsRequest,
) (hub.GetIncognitoPostCommentsResponse, error) {
	pg.log.Dbg("entered GetIncognitoPostComments", "id", req.IncognitoPostID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	// First check if the incognito post exists and is not deleted
	if !pg.incognitoPostExists(ctx, req.IncognitoPostID) {
		pg.log.Dbg("not found or deleted", "id", req.IncognitoPostID)
		return hub.GetIncognitoPostCommentsResponse{}, db.ErrNoIncognitoPost
	}

	// Single efficient query to get all comment data
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
			CASE WHEN ipc.is_deleted THEN TRUE ELSE FALSE END as is_deleted,
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
		pg.log.Err("failed to query incognito post comments",
			"error", err,
			"incognito_post_id", req.IncognitoPostID,
		)
		return hub.GetIncognitoPostCommentsResponse{}, err
	}
	defer rows.Close()

	var comments []hub.IncognitoPostComment
	var deletedCount int

	for rows.Next() {
		var comment hub.IncognitoPostComment
		var parentCommentID sql.NullString
		var isDeleted bool

		err := rows.Scan(
			&comment.CommentID,
			&comment.Content,
			&parentCommentID,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpvotesCount,
			&comment.DownvotesCount,
			&comment.IsCreatedByMe,
			&isDeleted,
			&comment.MeUpvoted,
			&comment.MeDownvoted,
		)
		if err != nil {
			pg.log.Err("failed to scan comment row", "error", err)
			return hub.GetIncognitoPostCommentsResponse{}, err
		}

		if parentCommentID.Valid {
			comment.InReplyTo = &parentCommentID.String
		}

		comment.IsDeleted = isDeleted
		if isDeleted {
			comment.Content = ""
			deletedCount++
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating comment rows", "error", err)
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	pg.log.Dbg("fetched incognito post comments",
		"incognito_post_id", req.IncognitoPostID,
		"total_comments", len(comments),
		"deleted_comments", deletedCount)

	return hub.GetIncognitoPostCommentsResponse{
		Comments: comments,
	}, nil
}

// Helper function to check if incognito post exists and is not deleted
func (pg *PG) incognitoPostExists(
	ctx context.Context,
	incognitoPostID string,
) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM incognito_posts
		WHERE id = $1 AND is_deleted = FALSE
	`

	err := pg.pool.QueryRow(ctx, query, incognitoPostID).Scan(&count)
	return err == nil && count > 0
}
