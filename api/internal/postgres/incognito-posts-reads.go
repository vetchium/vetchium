package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
			COALESCE(ip.upvotes_count, 0) as upvotes_count,
			COALESCE(ip.downvotes_count, 0) as downvotes_count,
			COALESCE(ip.score, 0) as score,
			CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
					AND ipv.vote_value = $3
				) THEN TRUE
				ELSE FALSE
			END as me_upvoted,
						CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
					AND ipv.vote_value = $4
				) THEN TRUE
				ELSE FALSE
			END as me_downvoted,
			CASE
				WHEN ip.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_upvote,
			CASE
				WHEN ip.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_downvote,
			ip.is_deleted,
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
		GROUP BY ip.id, ip.content, ip.created_at, ip.author_id, ip.upvotes_count, ip.downvotes_count, ip.score, ip.is_deleted
	`

	var post hub.IncognitoPost
	var tagIDs []string
	var tagNames []string

	err = pg.pool.QueryRow(ctx, query, req.IncognitoPostID, hubUserID,
		db.UpvoteValue, db.DownvoteValue).Scan(
		&post.IncognitoPostID,
		&post.Content,
		&post.CreatedAt,
		&post.IsCreatedByMe,
		&post.UpvotesCount,
		&post.DownvotesCount,
		&post.Score,
		&post.MeUpvoted,
		&post.MeDownvoted,
		&post.CanUpvote,
		&post.CanDownvote,
		&post.IsDeleted,
		&tagIDs,
		&tagNames,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("incognito post not found", "id", req.IncognitoPostID)
			return hub.IncognitoPost{}, db.ErrNoIncognitoPost
		}
		pg.log.Err("failed to query incognito post",
			"error", err,
			"id", req.IncognitoPostID)
		return hub.IncognitoPost{}, err
	}

	// Vote fields are now directly scanned from the query

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
		pg.log.Dbg("failed to get hub user ID", "error", err)
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	// First check if the incognito post exists and is not deleted
	if !pg.incognitoPostExists(ctx, req.IncognitoPostID) {
		pg.log.Dbg("not found or deleted", "id", req.IncognitoPostID)
		return hub.GetIncognitoPostCommentsResponse{}, db.ErrNoIncognitoPost
	}

	// Get top-level comments with pagination and total count in one query
	topLevelComments, nextPaginationKey, totalCount, err := pg.getTopLevelCommentsWithCount(
		ctx,
		req.IncognitoPostID,
		hubUserID,
		req.SortBy,
		req.PaginationKey,
		req.Limit,
	)
	if err != nil {
		pg.log.Dbg("failed to get top-level comments", "error", err)
		return hub.GetIncognitoPostCommentsResponse{}, err
	}

	allComments := make([]hub.IncognitoPostComment, 0)
	allComments = append(allComments, topLevelComments...)

	// Get preview replies for all top-level comments in one query
	if req.RepliesPreviewCount > 0 && len(topLevelComments) > 0 {
		replies, err := pg.getBulkCommentRepliesPreview(
			ctx,
			topLevelComments,
			hubUserID,
			req.RepliesPreviewCount,
		)
		if err != nil {
			pg.log.Dbg("failed to get replies preview", "error", err)
			return hub.GetIncognitoPostCommentsResponse{}, err
		}
		allComments = append(allComments, replies...)
	}

	pg.log.Dbg("fetched incognito post comments",
		"incognito_post_id", req.IncognitoPostID,
		"top_level_comments", len(topLevelComments),
		"total_comments_returned", len(allComments),
		"total_top_level_count", totalCount)

	return hub.GetIncognitoPostCommentsResponse{
		Comments:           allComments,
		PaginationKey:      nextPaginationKey,
		TotalCommentsCount: totalCount,
	}, nil
}

// getTopLevelCommentsWithCount fetches top-level comments with pagination and total count
func (pg *PG) getTopLevelCommentsWithCount(
	ctx context.Context,
	incognitoPostID string,
	hubUserID string,
	sortBy hub.IncognitoPostCommentSortBy,
	paginationKey *string,
	limit int32,
) ([]hub.IncognitoPostComment, string, int32, error) {
	// First get total count in a simple query
	var totalCount int32
	countQuery := `
		SELECT COUNT(*)
		FROM incognito_post_comments
		WHERE incognito_post_id = $1 AND parent_comment_id IS NULL
	`
	err := pg.pool.QueryRow(ctx, countQuery, incognitoPostID).Scan(&totalCount)
	if err != nil {
		pg.log.Err("failed to get total top-level comments count", "error", err)
		return nil, "", 0, db.ErrInternal
	}

	// Then get the paginated comments
	var query string
	var args []interface{}

	baseQuery := `
		SELECT
			ipc.id,
			ipc.content,
			ipc.parent_comment_id,
			ipc.depth,
			ipc.created_at,
			ipc.upvotes_count,
			ipc.downvotes_count,
			COALESCE(ipc.score, 0) as score,
			CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv
					WHERE ipcv.comment_id = ipc.id
					AND ipcv.user_id = $2
					AND ipcv.vote_value = $3
				) THEN TRUE
				ELSE FALSE
			END as me_upvoted,
			CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv
					WHERE ipcv.comment_id = ipc.id
					AND ipcv.user_id = $2
					AND ipcv.vote_value = $4
				) THEN TRUE
				ELSE FALSE
			END as me_downvoted,
			CASE
				WHEN ipc.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv
					WHERE ipcv.comment_id = ipc.id
					AND ipcv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_upvote,
			CASE
				WHEN ipc.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_comment_votes ipcv
					WHERE ipcv.comment_id = ipc.id
					AND ipcv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_downvote,
			CASE WHEN ipc.author_id = $2 THEN TRUE ELSE FALSE END as is_created_by_me,
			CASE WHEN ipc.is_deleted THEN TRUE ELSE FALSE END as is_deleted,
			COALESCE(
				(SELECT COUNT(*) FROM incognito_post_comments child
				 WHERE child.parent_comment_id = ipc.id), 0
			) as replies_count
		FROM incognito_post_comments ipc
		WHERE ipc.incognito_post_id = $1 AND ipc.parent_comment_id IS NULL
	`

	args = []interface{}{
		incognitoPostID,
		hubUserID,
		db.UpvoteValue,
		db.DownvoteValue,
	}

	// Add pagination logic
	if paginationKey != nil && *paginationKey != "" {
		switch sortBy {
		case hub.IncognitoPostCommentSortByTop:
			var lastScore sql.NullInt32
			var lastCreatedAt sql.NullTime
			err := pg.pool.QueryRow(
				ctx,
				`
				SELECT COALESCE(score, 0), created_at
				FROM incognito_post_comments
				WHERE id = $1 AND incognito_post_id = $2`,
				*paginationKey,
				incognitoPostID,
			).Scan(&lastScore, &lastCreatedAt)

			if err != nil {
				pg.log.Err("pagination cursor data for top sorting",
					"error", err,
					"pagination_key", *paginationKey,
				)
				return nil, "", 0, db.ErrInvalidPaginationKey
			}

			if lastScore.Valid && lastCreatedAt.Valid {
				baseQuery += ` AND (
					ipc.score < $5 OR
					(ipc.score = $5 AND ipc.created_at < $6) OR
					(ipc.score = $5 AND ipc.created_at = $6 AND ipc.id > $7)
				)`
				args = append(
					args,
					lastScore.Int32,
					lastCreatedAt.Time,
					*paginationKey,
				)
			}

		case hub.IncognitoPostCommentSortByNew:
			var lastCreatedAt sql.NullTime
			err := pg.pool.QueryRow(ctx, `
				SELECT created_at
				FROM incognito_post_comments
				WHERE id = $1 AND incognito_post_id = $2`,
				*paginationKey, incognitoPostID).Scan(&lastCreatedAt)

			if err != nil {
				pg.log.Err("pagination cursor data for new sorting",
					"error", err,
					"pagination_key", *paginationKey)
				return nil, "", 0, db.ErrInvalidPaginationKey
			}

			if lastCreatedAt.Valid {
				baseQuery += ` AND (
					ipc.created_at < $5 OR
					(ipc.created_at = $5 AND ipc.id > $6)
				)`
				args = append(args, lastCreatedAt.Time, *paginationKey)
			}

		case hub.IncognitoPostCommentSortByOld:
			var lastCreatedAt sql.NullTime
			err := pg.pool.QueryRow(ctx, `
				SELECT created_at
				FROM incognito_post_comments
				WHERE id = $1 AND incognito_post_id = $2`,
				*paginationKey, incognitoPostID).Scan(&lastCreatedAt)

			if err != nil {
				pg.log.Err("pagination cursor data for old sorting",
					"error", err,
					"pagination_key", *paginationKey)
				return nil, "", 0, db.ErrInvalidPaginationKey
			}

			if lastCreatedAt.Valid {
				baseQuery += ` AND (
					ipc.created_at > $5 OR
					(ipc.created_at = $5 AND ipc.id > $6)
				)`
				args = append(args, lastCreatedAt.Time, *paginationKey)
			}
		}
	}

	// Add ordering based on sort_by
	switch sortBy {
	case hub.IncognitoPostCommentSortByTop:
		baseQuery += ` ORDER BY ipc.score DESC, ipc.created_at DESC, ipc.id ASC`
	case hub.IncognitoPostCommentSortByNew:
		baseQuery += ` ORDER BY ipc.created_at DESC, ipc.id ASC`
	case hub.IncognitoPostCommentSortByOld:
		baseQuery += ` ORDER BY ipc.created_at ASC, ipc.id ASC`
	default:
		baseQuery += ` ORDER BY ipc.score DESC, ipc.created_at DESC, ipc.id ASC`
	}

	// Add limit
	baseQuery += fmt.Sprintf(` LIMIT $%d`, len(args)+1)
	args = append(args, limit)

	query = baseQuery

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed to query top-level comments", "error", err)
		return nil, "", 0, db.ErrInternal
	}
	defer rows.Close()

	comments := make([]hub.IncognitoPostComment, 0)
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
			&comment.Score,
			&comment.MeUpvoted,
			&comment.MeDownvoted,
			&comment.CanUpvote,
			&comment.CanDownvote,
			&comment.IsCreatedByMe,
			&isDeleted,
			&comment.RepliesCount,
		)
		if err != nil {
			pg.log.Err("failed to scan top-level comment row", "error", err)
			return nil, "", 0, db.ErrInternal
		}

		if parentCommentID.Valid {
			comment.InReplyTo = &parentCommentID.String
		}

		comment.IsDeleted = isDeleted
		if isDeleted {
			comment.Content = ""
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating top-level comment rows", "error", err)
		return nil, "", 0, db.ErrInternal
	}

	var nextPaginationKey string
	if len(comments) == int(limit) {
		nextPaginationKey = comments[len(comments)-1].CommentID
	}

	return comments, nextPaginationKey, totalCount, nil
}

// getBulkCommentRepliesPreview fetches preview replies for multiple comments in one query
func (pg *PG) getBulkCommentRepliesPreview(
	ctx context.Context,
	parentComments []hub.IncognitoPostComment,
	hubUserID string,
	limitPerComment int32,
) ([]hub.IncognitoPostComment, error) {
	if len(parentComments) == 0 {
		return []hub.IncognitoPostComment{}, nil
	}

	// Build array of parent comment IDs
	parentIDs := make([]string, len(parentComments))
	for i, comment := range parentComments {
		parentIDs[i] = comment.CommentID
	}

	query := `
		WITH ranked_replies AS (
			SELECT
				ipc.id,
				ipc.content,
				ipc.parent_comment_id,
				ipc.depth,
				ipc.created_at,
				ipc.upvotes_count,
				ipc.downvotes_count,
				COALESCE(ipc.score, 0) as score,
				CASE
					WHEN EXISTS(
						SELECT 1 FROM incognito_post_comment_votes ipcv
						WHERE ipcv.comment_id = ipc.id
						AND ipcv.user_id = $2
						AND ipcv.vote_value = $3
					) THEN TRUE
					ELSE FALSE
				END as me_upvoted,
				CASE
					WHEN EXISTS(
						SELECT 1 FROM incognito_post_comment_votes ipcv
						WHERE ipcv.comment_id = ipc.id
						AND ipcv.user_id = $2
						AND ipcv.vote_value = $4
					) THEN TRUE
					ELSE FALSE
				END as me_downvoted,
				CASE
					WHEN ipc.author_id = $2 THEN FALSE
					WHEN EXISTS(
						SELECT 1 FROM incognito_post_comment_votes ipcv
						WHERE ipcv.comment_id = ipc.id
						AND ipcv.user_id = $2
					) THEN FALSE
					ELSE TRUE
				END as can_upvote,
				CASE
					WHEN ipc.author_id = $2 THEN FALSE
					WHEN EXISTS(
						SELECT 1 FROM incognito_post_comment_votes ipcv
						WHERE ipcv.comment_id = ipc.id
						AND ipcv.user_id = $2
					) THEN FALSE
					ELSE TRUE
				END as can_downvote,
				CASE WHEN ipc.author_id = $2 THEN TRUE ELSE FALSE END as is_created_by_me,
				CASE WHEN ipc.is_deleted THEN TRUE ELSE FALSE END as is_deleted,
				COALESCE(
					(SELECT COUNT(*) FROM incognito_post_comments child
					 WHERE child.parent_comment_id = ipc.id), 0
				) as replies_count,
				ROW_NUMBER() OVER (
					PARTITION BY ipc.parent_comment_id
					ORDER BY ipc.score DESC, ipc.created_at ASC
				) as rn
			FROM incognito_post_comments ipc
			WHERE ipc.parent_comment_id = ANY($1)
		)
		SELECT
			id, content, parent_comment_id, depth, created_at,
			upvotes_count, downvotes_count, score,
			me_upvoted, me_downvoted, can_upvote, can_downvote,
			is_created_by_me, is_deleted, replies_count
		FROM ranked_replies
		WHERE rn <= $5
		ORDER BY parent_comment_id, rn
	`

	rows, err := pg.pool.Query(ctx, query, parentIDs, hubUserID,
		db.UpvoteValue, db.DownvoteValue, limitPerComment)
	if err != nil {
		pg.log.Err("failed to query bulk comment replies preview", "error", err)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	replies := make([]hub.IncognitoPostComment, 0)
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
			&comment.Score,
			&comment.MeUpvoted,
			&comment.MeDownvoted,
			&comment.CanUpvote,
			&comment.CanDownvote,
			&comment.IsCreatedByMe,
			&isDeleted,
			&comment.RepliesCount,
		)
		if err != nil {
			pg.log.Err("failed to scan reply comment row", "error", err)
			return nil, db.ErrInternal
		}

		if parentCommentID.Valid {
			comment.InReplyTo = &parentCommentID.String
		}

		comment.IsDeleted = isDeleted
		if isDeleted {
			comment.Content = ""
		}

		replies = append(replies, comment)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating reply comment rows", "error", err)
		return nil, db.ErrInternal
	}

	return replies, nil
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
	if err != nil {
		pg.log.Err("failed to check if incognito post exists", "error", err)
		return false
	}
	return count > 0
}

func (pg *PG) GetIncognitoPosts(
	ctx context.Context,
	req hub.GetIncognitoPostsRequest,
) (hub.GetIncognitoPostsResponse, error) {
	pg.log.Dbg("entered GetIncognitoPosts",
		"tag_id", req.TagID,
		"limit", req.Limit)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.GetIncognitoPostsResponse{}, err
	}

	// Build the base query with pagination and filtering
	var query string
	var args []interface{}

	baseQuery := `
		SELECT
			ip.id,
			ip.content,
			ip.created_at,
			CASE WHEN ip.author_id = $2 THEN TRUE ELSE FALSE END as is_created_by_me,
			COALESCE(ip.upvotes_count, 0) as upvotes_count,
			COALESCE(ip.downvotes_count, 0) as downvotes_count,
			COALESCE(ip.score, 0) as score,
			CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
					AND ipv.vote_value = $3
				) THEN TRUE
				ELSE FALSE
			END as me_upvoted,
			CASE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
					AND ipv.vote_value = $4
				) THEN TRUE
				ELSE FALSE
			END as me_downvoted,
			CASE
				WHEN ip.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_upvote,
			CASE
				WHEN ip.author_id = $2 THEN FALSE
				WHEN EXISTS(
					SELECT 1 FROM incognito_post_votes ipv
					WHERE ipv.incognito_post_id = ip.id
					AND ipv.user_id = $2
				) THEN FALSE
				ELSE TRUE
			END as can_downvote,
			ip.is_deleted,
			COALESCE(
				(SELECT COUNT(*) FROM incognito_post_comments ipc
				 WHERE ipc.incognito_post_id = ip.id AND ipc.is_deleted = FALSE), 0
			) as comments_count,
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
		WHERE ip.is_deleted = FALSE
		AND EXISTS (
			SELECT 1 FROM incognito_post_tags ipt2
			WHERE ipt2.incognito_post_id = ip.id AND ipt2.tag_id = $1
		)
	`

	args = []interface{}{req.TagID, hubUserID, db.UpvoteValue, db.DownvoteValue}

	// Add pagination logic
	if req.PaginationKey != nil && *req.PaginationKey != "" {
		// Get score and created_at of the pagination key post for cursor-based pagination
		var lastScore sql.NullInt32
		var lastCreatedAt sql.NullTime
		err := pg.pool.QueryRow(ctx, `
			SELECT COALESCE(score, 0), created_at
			FROM incognito_posts
			WHERE id = $1`,
			*req.PaginationKey).Scan(&lastScore, &lastCreatedAt)

		if err != nil {
			pg.log.Err("failed to get pagination cursor data", "error", err)
			return hub.GetIncognitoPostsResponse{}, err
		}

		if lastScore.Valid && lastCreatedAt.Valid {
			// Add pagination condition: posts with lower score OR same score but
			// older created_at OR same score and created_at but smaller ID
			baseQuery += ` AND (
				ip.score < $5 OR
				(ip.score = $5 AND ip.created_at < $6) OR
				(ip.score = $5 AND ip.created_at = $6 AND ip.id < $7)
			)`
			args = append(
				args,
				lastScore.Int32,
				lastCreatedAt.Time,
				*req.PaginationKey,
			)
		}
	}

	query = baseQuery + `
		GROUP BY ip.id, ip.content, ip.created_at, ip.author_id, ip.upvotes_count, ip.downvotes_count, ip.score, ip.is_deleted
		ORDER BY ip.score DESC, ip.created_at DESC, ip.id DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, req.Limit)

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed to query incognito posts",
			"error", err,
			"tag_id", req.TagID)
		return hub.GetIncognitoPostsResponse{}, err
	}
	defer rows.Close()

	posts := make([]hub.IncognitoPostSummary, 0)
	for rows.Next() {
		var post hub.IncognitoPostSummary
		var tagIDs []string
		var tagNames []string

		err := rows.Scan(
			&post.IncognitoPostID,
			&post.Content,
			&post.CreatedAt,
			&post.IsCreatedByMe,
			&post.UpvotesCount,
			&post.DownvotesCount,
			&post.Score,
			&post.MeUpvoted,
			&post.MeDownvoted,
			&post.CanUpvote,
			&post.CanDownvote,
			&post.IsDeleted,
			&post.CommentsCount,
			&tagIDs,
			&tagNames,
		)
		if err != nil {
			pg.log.Err("failed to scan post row", "error", err)
			return hub.GetIncognitoPostsResponse{}, err
		}

		// Build tags array
		post.Tags = make([]common.VTag, len(tagIDs))
		for i := 0; i < len(tagIDs) && i < len(tagNames); i++ {
			post.Tags[i] = common.VTag{
				ID:   common.VTagID(tagIDs[i]),
				Name: common.VTagName(tagNames[i]),
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating post rows", "error", err)
		return hub.GetIncognitoPostsResponse{}, err
	}

	pg.log.Dbg("fetched incognito posts",
		"tag_id", req.TagID,
		"count", len(posts))

	// Set pagination key if we reached the limit (indicating more data available)
	var paginationKey string
	if len(posts) == int(req.Limit) {
		// Use the last post's ID as pagination key
		paginationKey = posts[len(posts)-1].IncognitoPostID
	}

	return hub.GetIncognitoPostsResponse{
		Posts:         posts,
		PaginationKey: paginationKey,
	}, nil
}

func (pg *PG) GetMyIncognitoPosts(
	ctx context.Context,
	req hub.GetMyIncognitoPostsRequest,
) (hub.GetMyIncognitoPostsResponse, error) {
	pg.log.Dbg("entered GetMyIncognitoPosts", "limit", req.Limit)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.GetMyIncognitoPostsResponse{}, err
	}

	var query string
	var args []interface{}

	baseQuery := `
		SELECT
			ip.id,
			ip.content,
			ip.created_at,
			TRUE as is_created_by_me,
			COALESCE(ip.upvotes_count, 0) as upvotes_count,
			COALESCE(ip.downvotes_count, 0) as downvotes_count,
			COALESCE(ip.score, 0) as score,
			FALSE as me_upvoted,
			FALSE as me_downvoted,
			FALSE as can_upvote,
			FALSE as can_downvote,
			ip.is_deleted,
			COALESCE(
				(SELECT COUNT(*) FROM incognito_post_comments ipc
				 WHERE ipc.incognito_post_id = ip.id AND ipc.is_deleted = FALSE), 0
			) as comments_count,
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
		WHERE ip.author_id = $1`

	args = []interface{}{hubUserID}

	// Add pagination logic
	if req.PaginationKey != nil && *req.PaginationKey != "" {
		// Get created_at of the pagination key post for cursor-based pagination
		var lastCreatedAt sql.NullTime
		err := pg.pool.QueryRow(ctx, `
			SELECT created_at
			FROM incognito_posts
			WHERE id = $1 AND author_id = $2`,
			*req.PaginationKey, hubUserID).Scan(&lastCreatedAt)

		if err != nil {
			pg.log.Err("failed to get my posts pagination cursor data",
				"error", err)
			return hub.GetMyIncognitoPostsResponse{}, err
		}

		if lastCreatedAt.Valid {
			// Add pagination condition: posts older than the cursor OR same created_at but smaller ID
			baseQuery += ` AND (
				ip.created_at < $2 OR
				(ip.created_at = $2 AND ip.id < $3)
			)`
			args = append(args, lastCreatedAt.Time, *req.PaginationKey)
		}
	}

	query = baseQuery + `
		GROUP BY ip.id, ip.content, ip.created_at, ip.author_id, ip.upvotes_count, ip.downvotes_count, ip.score, ip.is_deleted
		ORDER BY ip.created_at DESC, ip.id DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, req.Limit)

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed to query my incognito posts", "error", err)
		return hub.GetMyIncognitoPostsResponse{}, err
	}
	defer rows.Close()

	posts := make([]hub.IncognitoPostSummary, 0)
	for rows.Next() {
		var post hub.IncognitoPostSummary
		var tagIDs []string
		var tagNames []string

		err := rows.Scan(
			&post.IncognitoPostID,
			&post.Content,
			&post.CreatedAt,
			&post.IsCreatedByMe,
			&post.UpvotesCount,
			&post.DownvotesCount,
			&post.Score,
			&post.MeUpvoted,
			&post.MeDownvoted,
			&post.CanUpvote,
			&post.CanDownvote,
			&post.IsDeleted,
			&post.CommentsCount,
			&tagIDs,
			&tagNames,
		)
		if err != nil {
			pg.log.Err("failed to scan my post row", "error", err)
			return hub.GetMyIncognitoPostsResponse{}, err
		}

		// Build tags array
		post.Tags = make([]common.VTag, len(tagIDs))
		for i := 0; i < len(tagIDs) && i < len(tagNames); i++ {
			post.Tags[i] = common.VTag{
				ID:   common.VTagID(tagIDs[i]),
				Name: common.VTagName(tagNames[i]),
			}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating my post rows", "error", err)
		return hub.GetMyIncognitoPostsResponse{}, err
	}

	pg.log.Dbg("fetched my incognito posts", "count", len(posts))

	// Set pagination key if we reached the limit (indicating more data available)
	var paginationKey string
	if len(posts) == int(req.Limit) {
		// Use the last post's ID as pagination key
		paginationKey = posts[len(posts)-1].IncognitoPostID
	}

	return hub.GetMyIncognitoPostsResponse{
		Posts:         posts,
		PaginationKey: paginationKey,
	}, nil
}

func (pg *PG) GetMyIncognitoPostComments(
	ctx context.Context,
	req hub.GetMyIncognitoPostCommentsRequest,
) (hub.GetMyIncognitoPostCommentsResponse, error) {
	pg.log.Dbg("entered GetMyIncognitoPostComments", "limit", req.Limit)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.GetMyIncognitoPostCommentsResponse{}, err
	}

	var query string
	var args []interface{}

	baseQuery := `
		SELECT
			ipc.id,
			ipc.incognito_post_id,
			ipc.content,
			ipc.parent_comment_id,
			ipc.depth,
			ipc.created_at,
			ipc.upvotes_count,
			ipc.downvotes_count,
			COALESCE(ipc.score, 0) as score,
			FALSE as me_upvoted,
			FALSE as me_downvoted,
			CASE WHEN ipc.is_deleted THEN TRUE ELSE FALSE END as is_deleted,
			LEFT(ip.content, 100) as post_content_preview,
			COALESCE(
				ARRAY_AGG(t.id ORDER BY t.display_name) FILTER (WHERE t.id IS NOT NULL),
				'{}'::text[]
			) as post_tag_ids,
			COALESCE(
				ARRAY_AGG(t.display_name ORDER BY t.display_name) FILTER (WHERE t.display_name IS NOT NULL),
				'{}'::text[]
			) as post_tag_names
		FROM incognito_post_comments ipc
		JOIN incognito_posts ip ON ipc.incognito_post_id = ip.id
		LEFT JOIN incognito_post_tags ipt ON ip.id = ipt.incognito_post_id
		LEFT JOIN tags t ON ipt.tag_id = t.id
		WHERE ipc.author_id = $1 AND ip.is_deleted = FALSE`

	args = []interface{}{hubUserID}

	// Add pagination logic
	if req.PaginationKey != nil && *req.PaginationKey != "" {
		// Get created_at of the pagination key comment for cursor-based pagination
		var lastCreatedAt sql.NullTime
		err := pg.pool.QueryRow(ctx, `
			SELECT created_at
			FROM incognito_post_comments
			WHERE id = $1 AND author_id = $2`,
			*req.PaginationKey, hubUserID).Scan(&lastCreatedAt)

		if err != nil {
			pg.log.Err("failed to get my comments pagination cursor data",
				"error", err)
			return hub.GetMyIncognitoPostCommentsResponse{}, err
		}

		if lastCreatedAt.Valid {
			// Add pagination condition: comments older than the cursor OR same
			// created_at but smaller ID
			baseQuery += ` AND (
				ipc.created_at < $2 OR
				(ipc.created_at = $2 AND ipc.id < $3)
			)`
			args = append(args, lastCreatedAt.Time, *req.PaginationKey)
		}
	}

	query = baseQuery + `
		GROUP BY ipc.id, ipc.incognito_post_id, ipc.content, ipc.parent_comment_id, ipc.depth, ipc.created_at, ipc.upvotes_count, ipc.downvotes_count, ipc.score, ipc.is_deleted, ip.content
		ORDER BY ipc.created_at DESC, ipc.id DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)+1)

	args = append(args, req.Limit)

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed to query my incognito post comments", "error", err)
		return hub.GetMyIncognitoPostCommentsResponse{}, err
	}
	defer rows.Close()

	comments := make([]hub.MyIncognitoPostComment, 0)
	for rows.Next() {
		var comment hub.MyIncognitoPostComment
		var parentCommentID sql.NullString
		var isDeleted bool
		var postTagIDs []string
		var postTagNames []string

		err := rows.Scan(
			&comment.CommentID,
			&comment.IncognitoPostID,
			&comment.Content,
			&parentCommentID,
			&comment.Depth,
			&comment.CreatedAt,
			&comment.UpvotesCount,
			&comment.DownvotesCount,
			&comment.Score,
			&comment.MeUpvoted,
			&comment.MeDownvoted,
			&isDeleted,
			&comment.PostContentPreview,
			&postTagIDs,
			&postTagNames,
		)
		if err != nil {
			pg.log.Err("failed to scan my comment row", "error", err)
			return hub.GetMyIncognitoPostCommentsResponse{}, err
		}

		if parentCommentID.Valid {
			comment.InReplyTo = &parentCommentID.String
		}

		comment.IsDeleted = isDeleted
		if isDeleted {
			comment.Content = ""
		}

		// Build post tags array
		comment.PostTags = make([]common.VTag, len(postTagIDs))
		for i := 0; i < len(postTagIDs) && i < len(postTagNames); i++ {
			comment.PostTags[i] = common.VTag{
				ID:   common.VTagID(postTagIDs[i]),
				Name: common.VTagName(postTagNames[i]),
			}
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating my comment rows", "error", err)
		return hub.GetMyIncognitoPostCommentsResponse{}, err
	}

	pg.log.Dbg("fetched my incognito post comments", "count", len(comments))

	// Set pagination key if we reached the limit (indicating more data available)
	var paginationKey string
	if len(comments) == int(req.Limit) {
		// Use the last comment's ID as pagination key
		paginationKey = comments[len(comments)-1].CommentID
	}

	return hub.GetMyIncognitoPostCommentsResponse{
		Comments:      comments,
		PaginationKey: paginationKey,
	}, nil
}
