package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) AddIncognitoPost(
	ctx context.Context,
	req db.AddIncognitoPostRequest,
) error {
	pg.log.Dbg(
		"entered AddIncognitoPost",
		"incognito_post_id",
		req.IncognitoPostID,
	)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		pg.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Validate tag IDs if provided
	if len(req.AddIncognitoPostReq.TagIDs) > 0 {
		tagIDStrings := make([]string, len(req.AddIncognitoPostReq.TagIDs))
		for i, tagID := range req.AddIncognitoPostReq.TagIDs {
			tagIDStrings[i] = string(tagID)
		}

		if !pg.validateTagIDs(ctx, tagIDStrings) {
			pg.log.Dbg("invalid tag IDs provided", "tag_ids", tagIDStrings)
			return db.ErrInvalidTagIDs
		}
	}

	// Insert the post
	query := `
		INSERT INTO incognito_posts (id, content, author_id)
		VALUES ($1, $2, $3)
	`

	_, err = tx.Exec(
		ctx,
		query,
		req.IncognitoPostID,
		req.AddIncognitoPostReq.Content,
		hubUserID,
	)
	if err != nil {
		pg.log.Err(
			"failed to insert incognito post",
			"error",
			err,
			"incognito_post_id",
			req.IncognitoPostID,
		)
		return err
	}

	// Insert tags if provided
	if len(req.AddIncognitoPostReq.TagIDs) > 0 {
		tagIDStrings := make([]string, len(req.AddIncognitoPostReq.TagIDs))
		for i, tagID := range req.AddIncognitoPostReq.TagIDs {
			tagIDStrings[i] = string(tagID)
		}

		for _, tagID := range tagIDStrings {
			tagQuery := `
				INSERT INTO incognito_post_tags (incognito_post_id, tag_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`
			_, err = tx.Exec(ctx, tagQuery, req.IncognitoPostID, tagID)
			if err != nil {
				pg.log.Err(
					"failed to insert incognito post tag",
					"error",
					err,
					"tag_id",
					tagID,
				)
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return err
	}

	pg.log.Dbg("created incognito post",
		"incognito_post_id", req.IncognitoPostID,
		"tag_count", len(req.AddIncognitoPostReq.TagIDs))
	return nil
}

func (pg *PG) DeleteIncognitoPost(
	ctx context.Context,
	req hub.DeleteIncognitoPostRequest,
) error {
	pg.log.Dbg(
		"entered DeleteIncognitoPost",
		"incognito_post_id",
		req.IncognitoPostID,
	)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// First check if the post exists and is not deleted
	var authorID string
	checkQuery := `
		SELECT author_id 
		FROM incognito_posts 
		WHERE id = $1 AND is_deleted = FALSE
	`

	err = pg.pool.QueryRow(ctx, checkQuery, req.IncognitoPostID).Scan(&authorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			pg.log.Dbg("incognito post not found",
				"incognito_post_id", req.IncognitoPostID)
			return db.ErrNoIncognitoPost
		}
		pg.log.Err("failed to check incognito post",
			"error", err,
			"incognito_post_id", req.IncognitoPostID)
		return err
	}

	// Check if the user is the author
	if authorID != hubUserID {
		pg.log.Dbg("user is not the author of the incognito post",
			"incognito_post_id", req.IncognitoPostID,
			"hub_user_id", hubUserID,
			"author_id", authorID)
		return db.ErrNotIncognitoPostAuthor
	}

	// Soft delete - set is_deleted = true
	updateQuery := `
		UPDATE incognito_posts
		SET is_deleted = TRUE
		WHERE id = $1 AND author_id = $2 AND is_deleted = FALSE
	`

	_, err = pg.pool.Exec(ctx, updateQuery, req.IncognitoPostID, hubUserID)
	if err != nil {
		pg.log.Err(
			"failed to soft delete incognito post",
			"error",
			err,
			"incognito_post_id",
			req.IncognitoPostID,
		)
		return err
	}

	pg.log.Dbg(
		"soft deleted incognito post",
		"incognito_post_id",
		req.IncognitoPostID,
	)
	return nil
}

func (pg *PG) AddIncognitoPostComment(
	ctx context.Context,
	req db.AddIncognitoPostCommentRequest,
) (hub.AddIncognitoPostCommentResponse, error) {
	pg.log.Dbg("entered AddIncognitoPostComment",
		"incognito_post_id", req.AddIncognitoPostCommentRequest.IncognitoPostID,
		"comment_id", req.CommentID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.AddIncognitoPostCommentResponse{}, err
	}

	// Check if incognito post exists and is not deleted
	if !pg.incognitoPostExists(
		ctx,
		req.AddIncognitoPostCommentRequest.IncognitoPostID,
	) {
		pg.log.Dbg(
			"incognito post not found or deleted",
			"incognito_post_id",
			req.AddIncognitoPostCommentRequest.IncognitoPostID,
		)
		return hub.AddIncognitoPostCommentResponse{}, db.ErrNoIncognitoPost
	}

	depth := 0
	if req.AddIncognitoPostCommentRequest.InReplyTo != nil {
		// Validate parent comment exists and is not deleted
		parentDepth, err := pg.getCommentDepth(
			ctx,
			*req.AddIncognitoPostCommentRequest.InReplyTo,
		)
		if err != nil {
			pg.log.Err(
				"failed to get parent comment depth",
				"error",
				err,
				"parent_comment_id",
				*req.AddIncognitoPostCommentRequest.InReplyTo,
			)
			return hub.AddIncognitoPostCommentResponse{}, err
		}

		// Check if adding a reply would exceed max depth
		if parentDepth >= 10 {
			pg.log.Dbg(
				"comment depth limit reached",
				"parent_depth",
				parentDepth,
				"parent_comment_id",
				*req.AddIncognitoPostCommentRequest.InReplyTo,
			)
			return hub.AddIncognitoPostCommentResponse{}, db.ErrMaxCommentDepthReached
		}

		depth = parentDepth + 1
	}

	query := `
		INSERT INTO incognito_post_comments (
			id, incognito_post_id, author_id, content, parent_comment_id, depth
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = pg.pool.Exec(
		ctx,
		query,
		req.CommentID,
		req.AddIncognitoPostCommentRequest.IncognitoPostID,
		hubUserID,
		req.AddIncognitoPostCommentRequest.Content,
		req.AddIncognitoPostCommentRequest.InReplyTo,
		depth,
	)
	if err != nil {
		pg.log.Err(
			"failed to insert incognito post comment",
			"error",
			err,
			"comment_id",
			req.CommentID,
		)
		return hub.AddIncognitoPostCommentResponse{}, err
	}

	response := hub.AddIncognitoPostCommentResponse{
		IncognitoPostID: req.AddIncognitoPostCommentRequest.IncognitoPostID,
		CommentID:       req.CommentID,
	}

	pg.log.Dbg("created incognito post comment",
		"incognito_post_id", req.AddIncognitoPostCommentRequest.IncognitoPostID,
		"comment_id", req.CommentID,
		"depth", depth)
	return response, nil
}

func (pg *PG) DeleteIncognitoPostComment(
	ctx context.Context,
	req hub.DeleteIncognitoPostCommentRequest,
) error {
	pg.log.Dbg("entered DeleteIncognitoPostComment",
		"incognito_post_id", req.IncognitoPostID,
		"comment_id", req.CommentID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Check if incognito post exists
	if !pg.incognitoPostExists(ctx, req.IncognitoPostID) {
		pg.log.Dbg(
			"incognito post not found",
			"incognito_post_id",
			req.IncognitoPostID,
		)
		return db.ErrNoIncognitoPost
	}

	// First check if the comment exists and is not deleted
	var authorID string
	checkQuery := `
		SELECT author_id 
		FROM incognito_post_comments 
		WHERE id = $1 AND incognito_post_id = $2 AND is_deleted = FALSE
	`

	err = pg.pool.QueryRow(ctx, checkQuery, req.CommentID, req.IncognitoPostID).
		Scan(&authorID)
	if err != nil {
		if err == pgx.ErrNoRows {
			pg.log.Dbg("incognito post comment not found",
				"comment_id", req.CommentID,
				"incognito_post_id", req.IncognitoPostID)
			return db.ErrNoIncognitoPostComment
		}
		pg.log.Err("failed to check incognito post comment",
			"error", err,
			"comment_id", req.CommentID)
		return err
	}

	// Check if the user is the author
	if authorID != hubUserID {
		pg.log.Dbg("user is not the author of the incognito post comment",
			"comment_id", req.CommentID,
			"incognito_post_id", req.IncognitoPostID,
			"hub_user_id", hubUserID,
			"author_id", authorID)
		return db.ErrNotIncognitoPostCommentAuthor
	}

	// Soft delete comment - set is_deleted = true
	updateQuery := `
		UPDATE incognito_post_comments
		SET is_deleted = TRUE
		WHERE id = $1 
		AND incognito_post_id = $2 
		AND author_id = $3
		AND is_deleted = FALSE
	`

	_, err = pg.pool.Exec(
		ctx,
		updateQuery,
		req.CommentID,
		req.IncognitoPostID,
		hubUserID,
	)
	if err != nil {
		pg.log.Err(
			"failed to soft delete incognito post comment",
			"error",
			err,
			"comment_id",
			req.CommentID,
		)
		return err
	}

	pg.log.Dbg("soft deleted incognito post comment",
		"incognito_post_id", req.IncognitoPostID,
		"comment_id", req.CommentID)
	return nil
}

func (pg *PG) UpvoteIncognitoPostComment(
	ctx context.Context,
	req hub.UpvoteIncognitoPostCommentRequest,
) error {
	pg.log.Dbg("entered UpvoteIncognitoPostComment",
		"incognito_post_id", req.IncognitoPostID,
		"comment_id", req.CommentID)

	err := pg.voteIncognitoPostComment(
		ctx,
		req.IncognitoPostID,
		req.CommentID,
		1,
	)
	if err != nil {
		pg.log.Err(
			"failed to upvote incognito post comment",
			"error",
			err,
			"comment_id",
			req.CommentID,
		)
		return err
	}

	pg.log.Dbg("upvoted incognito post comment", "comment_id", req.CommentID)
	return nil
}

func (pg *PG) DownvoteIncognitoPostComment(
	ctx context.Context,
	req hub.DownvoteIncognitoPostCommentRequest,
) error {
	pg.log.Dbg("entered DownvoteIncognitoPostComment",
		"incognito_post_id", req.IncognitoPostID,
		"comment_id", req.CommentID)

	err := pg.voteIncognitoPostComment(
		ctx,
		req.IncognitoPostID,
		req.CommentID,
		-1,
	)
	if err != nil {
		pg.log.Err(
			"failed to downvote incognito post comment",
			"error",
			err,
			"comment_id",
			req.CommentID,
		)
		return err
	}

	pg.log.Dbg("downvoted incognito post comment", "comment_id", req.CommentID)
	return nil
}

func (pg *PG) UnvoteIncognitoPostComment(
	ctx context.Context,
	req hub.UnvoteIncognitoPostCommentRequest,
) error {
	pg.log.Dbg("entered UnvoteIncognitoPostComment",
		"incognito_post_id", req.IncognitoPostID,
		"comment_id", req.CommentID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Single query to check comment exists, not deleted, and remove vote
	query := `
		WITH comment_check AS (
			SELECT id FROM incognito_post_comments 
			WHERE id = $1 AND incognito_post_id = $2 AND is_deleted = FALSE
		),
		vote_check AS (
			SELECT 1 FROM incognito_post_comment_votes ipcv
			JOIN comment_check cc ON TRUE
			WHERE ipcv.comment_id = $1 AND ipcv.user_id = $3
		)
		DELETE FROM incognito_post_comment_votes
		WHERE comment_id = $1 AND user_id = $3
		AND EXISTS (SELECT 1 FROM comment_check)
		AND EXISTS (SELECT 1 FROM vote_check)
	`

	result, err := pg.pool.Exec(
		ctx,
		query,
		req.CommentID,
		req.IncognitoPostID,
		hubUserID,
	)
	if err != nil {
		pg.log.Err(
			"failed to unvote incognito post comment",
			"error",
			err,
			"comment_id",
			req.CommentID,
		)
		return err
	}

	if result.RowsAffected() == 0 {
		// Check if comment exists but user can't vote (is author or comment doesn't exist)
		if !pg.commentExistsAndNotDeleted(
			ctx,
			req.CommentID,
			req.IncognitoPostID,
		) {
			pg.log.Dbg(
				"incognito post comment not found",
				"comment_id",
				req.CommentID,
			)
			return db.ErrNoIncognitoPostComment
		}
		if !pg.canVoteOnComment(ctx, req.CommentID, hubUserID) {
			pg.log.Dbg("user cannot vote on this comment",
				"comment_id", req.CommentID,
				"hub_user_id", hubUserID)
			return db.ErrNonVoteableIncognitoPostComment
		}
		pg.log.Dbg("no existing vote to remove",
			"comment_id", req.CommentID,
			"hub_user_id", hubUserID)
	}

	pg.log.Dbg("unvoted incognito post comment", "comment_id", req.CommentID)
	return nil
}

func (pg *PG) UnvoteIncognitoPost(
	ctx context.Context,
	req hub.UnvoteIncognitoPostRequest,
) error {
	pg.log.Dbg("entered UnvoteIncognitoPost", "id", req.IncognitoPostID)

	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Single query that checks post exists, user is not author, and removes vote
	query := `
		WITH post_check AS (
			SELECT author_id 
			FROM incognito_posts 
			WHERE id = $1 AND is_deleted = FALSE
		)
		DELETE FROM incognito_post_votes
		WHERE incognito_post_id = $1 AND user_id = $2
		AND EXISTS (
			SELECT 1 FROM post_check 
			WHERE author_id != $2
		)
	`

	result, err := pg.pool.Exec(ctx, query, req.IncognitoPostID, hubUserID)
	if err != nil {
		pg.log.Err("failed to unvote incognito post", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		// Check if post exists to determine the specific error
		var authorID string
		checkQuery := `
			SELECT author_id 
			FROM incognito_posts 
			WHERE id = $1 AND is_deleted = FALSE
		`
		err = pg.pool.QueryRow(ctx, checkQuery, req.IncognitoPostID).
			Scan(&authorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				pg.log.Dbg("incognitopost notfound", "id", req.IncognitoPostID)
				return db.ErrNoIncognitoPost
			}
			pg.log.Err("failed to check if post exists", "error", err)
			return err
		}

		// Post exists but user is the author
		if authorID == hubUserID {
			pg.log.Dbg("user is the author", "id", req.IncognitoPostID)
			return db.ErrNonVoteableIncognitoPost
		}

		// Post exists, user is not author, but no vote existed to remove (this is fine)
		pg.log.Dbg("no existing vote to remove",
			"incognito_post_id", req.IncognitoPostID,
			"hub_user_id", hubUserID,
		)
	}

	pg.log.Dbg("unvoted", "incognito_post_id", req.IncognitoPostID)
	return nil
}

func (pg *PG) voteIncognitoPostComment(
	ctx context.Context,
	incognitoPostID string,
	commentID string,
	voteValue int,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	// Check if incognito post exists
	if !pg.incognitoPostExists(ctx, incognitoPostID) {
		return db.ErrNoIncognitoPost
	}

	// Check if comment exists and is not deleted
	if !pg.commentExistsAndNotDeleted(ctx, commentID, incognitoPostID) {
		return db.ErrNoIncognitoPostComment
	}

	// Check if user can vote (not the author)
	if !pg.canVoteOnComment(ctx, commentID, hubUserID) {
		return db.ErrNonVoteableIncognitoPostComment
	}

	// Check for existing conflicting vote
	var existingVote *int
	voteCheckQuery := `
		SELECT vote_value 
		FROM incognito_post_comment_votes 
		WHERE comment_id = $1 AND user_id = $2
	`
	err = pg.pool.QueryRow(ctx, voteCheckQuery, commentID, hubUserID).
		Scan(&existingVote)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		pg.log.Err("failed to check existing comment vote", "error", err)
		return err
	}

	// If there's an existing vote and it's different from the new vote, return conflict error
	if existingVote != nil && *existingVote != voteValue {
		pg.log.Dbg("comment vote conflict detected",
			"comment_id", commentID,
			"existing_vote", *existingVote,
			"new_vote", voteValue)
		return db.ErrIncognitoPostCommentVoteConflict
	}

	// Insert or update vote (this will be idempotent for same vote)
	query := `
		INSERT INTO incognito_post_comment_votes (comment_id, user_id, vote_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (comment_id, user_id)
		DO UPDATE SET vote_value = EXCLUDED.vote_value
	`

	_, err = pg.pool.Exec(ctx, query, commentID, hubUserID, voteValue)
	return err
}

func (pg *PG) commentExistsAndNotDeleted(
	ctx context.Context,
	commentID, incognitoPostID string,
) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM incognito_post_comments
		WHERE id = $1 AND incognito_post_id = $2 AND is_deleted = FALSE
	`

	err := pg.pool.QueryRow(ctx, query, commentID, incognitoPostID).Scan(&count)
	return err == nil && count > 0
}

func (pg *PG) canVoteOnComment(
	ctx context.Context,
	commentID, hubUserID string,
) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM incognito_post_comments
		WHERE id = $1 AND author_id != $2 AND is_deleted = FALSE
	`

	err := pg.pool.QueryRow(ctx, query, commentID, hubUserID).Scan(&count)
	return err == nil && count > 0
}

// getCommentDepth returns the depth of a comment and validates it exists and is not deleted
func (pg *PG) getCommentDepth(
	ctx context.Context,
	commentID string,
) (int, error) {
	var depth int
	query := `
		SELECT depth
		FROM incognito_post_comments
		WHERE id = $1 AND is_deleted = FALSE
	`

	err := pg.pool.QueryRow(ctx, query, commentID).Scan(&depth)
	if err != nil {
		return 0, db.ErrInvalidParentComment
	}

	return depth, nil
}

func (pg *PG) validateTagIDs(ctx context.Context, tagIDs []string) bool {
	if len(tagIDs) == 0 {
		return true
	}

	// Deduplicate tag IDs
	uniqueTagIDs := make(map[string]bool)
	var deduplicatedTagIDs []string
	for _, tagID := range tagIDs {
		if !uniqueTagIDs[tagID] {
			uniqueTagIDs[tagID] = true
			deduplicatedTagIDs = append(deduplicatedTagIDs, tagID)
		}
	}

	query := `
		SELECT COUNT(DISTINCT id)
		FROM tags
		WHERE id = ANY($1)
	`

	var count int
	err := pg.pool.QueryRow(ctx, query, deduplicatedTagIDs).Scan(&count)
	return err == nil && count == len(deduplicatedTagIDs)
}

func (pg *PG) UpvoteIncognitoPost(
	ctx context.Context,
	req hub.UpvoteIncognitoPostRequest,
) error {
	pg.log.Dbg("entered UpvoteIncognitoPost",
		"incognito_post_id", req.IncognitoPostID)

	err := pg.voteIncognitoPost(ctx, req.IncognitoPostID, 1)
	if err != nil {
		pg.log.Dbg("failed to upvote incognito post",
			"error", err,
			"incognito_post_id", req.IncognitoPostID,
		)
		return err
	}

	pg.log.Dbg("upvoted", "incognito_post_id", req.IncognitoPostID)
	return nil
}

func (pg *PG) DownvoteIncognitoPost(
	ctx context.Context,
	req hub.DownvoteIncognitoPostRequest,
) error {
	pg.log.Dbg("entered DownvoteIncognitoPost", "id", req.IncognitoPostID)

	err := pg.voteIncognitoPost(ctx, req.IncognitoPostID, -1)
	if err != nil {
		pg.log.Dbg("failed to downvote incognito post", "error", err)
		return err
	}

	pg.log.Dbg("downvoted", "incognito_post_id", req.IncognitoPostID)
	return nil
}

func (pg *PG) voteIncognitoPost(
	ctx context.Context,
	incognitoPostID string,
	voteValue int,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	// Check if post exists and user is not the author
	var authorID string
	postCheckQuery := `
		SELECT author_id 
		FROM incognito_posts 
		WHERE id = $1 AND is_deleted = FALSE
	`
	err = pg.pool.QueryRow(ctx, postCheckQuery, incognitoPostID).Scan(&authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("incognito post not found", "id", incognitoPostID)
			return db.ErrNoIncognitoPost
		}
		pg.log.Err("failed to check if post exists", "error", err)
		return err
	}

	// Check if user is the author
	if authorID == hubUserID {
		pg.log.Dbg("user is the author", "id", incognitoPostID)
		return db.ErrNonVoteableIncognitoPost
	}

	// Check for existing conflicting vote
	var existingVote *int
	voteCheckQuery := `
		SELECT vote_value 
		FROM incognito_post_votes 
		WHERE incognito_post_id = $1 AND user_id = $2
	`
	err = pg.pool.QueryRow(ctx, voteCheckQuery, incognitoPostID, hubUserID).
		Scan(&existingVote)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		pg.log.Err("failed to check existing vote", "error", err)
		return err
	}

	// If there's an existing vote and it's different from the new vote, return conflict error
	if existingVote != nil && *existingVote != voteValue {
		pg.log.Dbg("vote conflict detected",
			"incognito_post_id", incognitoPostID,
			"existing_vote", *existingVote,
			"new_vote", voteValue)
		return db.ErrIncognitoPostVoteConflict
	}

	// Insert or update vote (this will be idempotent for same vote)
	insertQuery := `
		INSERT INTO incognito_post_votes (incognito_post_id, user_id, vote_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (incognito_post_id, user_id)
		DO UPDATE SET vote_value = EXCLUDED.vote_value
	`

	_, err = pg.pool.Exec(
		ctx,
		insertQuery,
		incognitoPostID,
		hubUserID,
		voteValue,
	)
	if err != nil {
		pg.log.Err("failed to vote incognito post", "error", err)
		return err
	}

	pg.log.Dbg("voted", "incognito_post_id", incognitoPostID)
	return nil
}
