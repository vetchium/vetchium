package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) AddIncognitoPost(
	ctx context.Context,
	req db.AddIncognitoPostRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

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
		return err
	}

	if len(req.AddIncognitoPostReq.TagIDs) > 0 {
		tagIDStrings := make([]string, len(req.AddIncognitoPostReq.TagIDs))
		for i, tagID := range req.AddIncognitoPostReq.TagIDs {
			tagIDStrings[i] = string(tagID)
		}
		err = pg.addIncognitoPostTags(
			ctx,
			tx,
			req.IncognitoPostID,
			tagIDStrings,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (pg *PG) DeleteIncognitoPost(
	ctx context.Context,
	req hub.DeleteIncognitoPostRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	query := `
		DELETE FROM incognito_posts
		WHERE id = $1 AND author_id = $2
	`

	result, err := pg.pool.Exec(ctx, query, req.IncognitoPostID, hubUserID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return db.ErrNoIncognitoPost
	}

	return nil
}

func (pg *PG) AddIncognitoPostComment(
	ctx context.Context,
	req db.AddIncognitoPostCommentRequest,
) (hub.AddIncognitoPostCommentResponse, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return hub.AddIncognitoPostCommentResponse{}, err
	}

	if !pg.canAccessIncognitoPost(
		ctx,
		req.AddIncognitoPostCommentRequest.IncognitoPostID,
		hubUserID,
	) {
		return hub.AddIncognitoPostCommentResponse{}, db.ErrNoIncognitoPost
	}

	depth := 0
	if req.AddIncognitoPostCommentRequest.InReplyTo != nil {
		calculatedDepth, err := pg.calculateCommentDepth(
			ctx,
			*req.AddIncognitoPostCommentRequest.InReplyTo,
		)
		if err != nil {
			return hub.AddIncognitoPostCommentResponse{}, err
		}
		depth = calculatedDepth
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
		return hub.AddIncognitoPostCommentResponse{}, err
	}

	return hub.AddIncognitoPostCommentResponse{
		IncognitoPostID: req.AddIncognitoPostCommentRequest.IncognitoPostID,
		CommentID:       req.CommentID,
	}, nil
}

func (pg *PG) DeleteIncognitoPostComment(
	ctx context.Context,
	req hub.DeleteIncognitoPostCommentRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	if !pg.canAccessIncognitoPost(ctx, req.IncognitoPostID, hubUserID) {
		return db.ErrNoIncognitoPost
	}

	query := `
		DELETE FROM incognito_post_comments
		WHERE id = $1 
		AND incognito_post_id = $2 
		AND author_id = $3
	`

	result, err := pg.pool.Exec(
		ctx,
		query,
		req.CommentID,
		req.IncognitoPostID,
		hubUserID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return db.ErrNoIncognitoPostComment
	}

	return nil
}

func (pg *PG) UpvoteIncognitoPostComment(
	ctx context.Context,
	req hub.UpvoteIncognitoPostCommentRequest,
) error {
	return pg.voteIncognitoPostComment(
		ctx,
		req.IncognitoPostID,
		req.CommentID,
		1,
	)
}

func (pg *PG) DownvoteIncognitoPostComment(
	ctx context.Context,
	req hub.DownvoteIncognitoPostCommentRequest,
) error {
	return pg.voteIncognitoPostComment(
		ctx,
		req.IncognitoPostID,
		req.CommentID,
		-1,
	)
}

func (pg *PG) UnvoteIncognitoPostComment(
	ctx context.Context,
	req hub.UnvoteIncognitoPostCommentRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		return err
	}

	if !pg.canAccessIncognitoPost(ctx, req.IncognitoPostID, hubUserID) {
		return db.ErrNoIncognitoPost
	}

	if !pg.commentExists(ctx, req.CommentID, req.IncognitoPostID) {
		return db.ErrNoIncognitoPostComment
	}

	if !pg.canVoteOnComment(ctx, req.CommentID, hubUserID) {
		return db.ErrNonVoteableIncognitoPostComment
	}

	query := `
		DELETE FROM incognito_post_comment_votes
		WHERE comment_id = $1 AND user_id = $2
	`

	_, err = pg.pool.Exec(ctx, query, req.CommentID, hubUserID)
	return err
}

// Helper functions

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

	if !pg.canAccessIncognitoPost(ctx, incognitoPostID, hubUserID) {
		return db.ErrNoIncognitoPost
	}

	if !pg.commentExists(ctx, commentID, incognitoPostID) {
		return db.ErrNoIncognitoPostComment
	}

	if !pg.canVoteOnComment(ctx, commentID, hubUserID) {
		return db.ErrNonVoteableIncognitoPostComment
	}

	query := `
		INSERT INTO incognito_post_comment_votes (comment_id, user_id, vote_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (comment_id, user_id)
		DO UPDATE SET vote_value = EXCLUDED.vote_value
	`

	_, err = pg.pool.Exec(ctx, query, commentID, hubUserID, voteValue)
	return err
}

func (pg *PG) commentExists(
	ctx context.Context,
	commentID, incognitoPostID string,
) bool {
	var count int
	query := `
		SELECT COUNT(*)
		FROM incognito_post_comments
		WHERE id = $1 AND incognito_post_id = $2
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
		WHERE id = $1 AND author_id != $2
	`

	err := pg.pool.QueryRow(ctx, query, commentID, hubUserID).Scan(&count)
	return err == nil && count > 0
}

func (pg *PG) calculateCommentDepth(
	ctx context.Context,
	parentCommentID string,
) (int, error) {
	var depth int
	query := `SELECT calculate_comment_depth($1)`

	err := pg.pool.QueryRow(ctx, query, parentCommentID).Scan(&depth)
	if err != nil {
		return 0, db.ErrInvalidParentComment
	}

	return depth, nil
}

func (pg *PG) addIncognitoPostTags(
	ctx context.Context,
	tx interface{},
	incognitoPostID string,
	tagIDs []string,
) error {
	// Check if all tag IDs exist
	if !pg.validateTagIDs(ctx, tagIDs) {
		return db.ErrInvalidTagIDs
	}

	for _, tagID := range tagIDs {
		query := `
			INSERT INTO incognito_post_tags (incognito_post_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, err := tx.(interface {
			Exec(ctx context.Context, sql string, arguments ...interface{}) (interface{}, error)
		}).Exec(ctx, query, incognitoPostID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pg *PG) validateTagIDs(ctx context.Context, tagIDs []string) bool {
	if len(tagIDs) == 0 {
		return true
	}

	query := `
		SELECT COUNT(DISTINCT id)
		FROM tags
		WHERE id = ANY($1)
	`

	var count int
	err := pg.pool.QueryRow(ctx, query, tagIDs).Scan(&count)
	return err == nil && count == len(tagIDs)
}
