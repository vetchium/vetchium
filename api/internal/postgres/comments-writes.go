package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) AddPostComment(
	ctx context.Context,
	req db.AddPostCommentRequest,
) (hub.AddPostCommentResponse, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return hub.AddPostCommentResponse{}, err
	}

	// Check if post exists and comments are enabled
	var postExists bool
	var commentsEnabled bool
	err = pg.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1),
		       COALESCE((SELECT comments_enabled FROM posts WHERE id = $1), false)
	`, req.PostID).Scan(&postExists, &commentsEnabled)
	if err != nil {
		pg.log.Err("failed to check post existence and comments status",
			"error", err, "post_id", req.PostID)
		return hub.AddPostCommentResponse{}, db.ErrInternal
	}

	if !postExists {
		pg.log.Dbg("post not found", "post_id", req.PostID)
		return hub.AddPostCommentResponse{}, db.ErrNoPost
	}

	if !commentsEnabled {
		pg.log.Dbg("comments disabled for post", "post_id", req.PostID)
		return hub.AddPostCommentResponse{}, db.ErrCommentsDisabled
	}

	// Insert comment
	_, err = pg.pool.Exec(ctx, `
		INSERT INTO post_comments (id, post_id, author_id, content)
		VALUES ($1, $2, $3, $4)
	`, req.CommentID, req.PostID, hubUserID, req.Content)
	if err != nil {
		pg.log.Err("failed to insert comment",
			"error", err, "comment_id", req.CommentID, "post_id", req.PostID)
		return hub.AddPostCommentResponse{}, db.ErrInternal
	}

	return hub.AddPostCommentResponse{
		PostID:    req.PostID,
		CommentID: req.CommentID,
	}, nil
}

func (pg *PG) DisablePostComments(
	ctx context.Context,
	req hub.DisablePostCommentsRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		pg.log.Err("failed to begin transaction", "error", err)
		return db.ErrInternal
	}
	defer tx.Rollback(context.Background())

	// Check if post exists and user is the author
	var postExists bool
	var isAuthor bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1),
		       EXISTS(SELECT 1 FROM posts WHERE id = $1 AND author_id = $2)
	`, req.PostID, hubUserID).Scan(&postExists, &isAuthor)
	if err != nil {
		pg.log.Err("failed to check post ownership",
			"error", err, "post_id", req.PostID, "user_id", hubUserID)
		return db.ErrInternal
	}

	if !postExists {
		pg.log.Dbg("post not found", "post_id", req.PostID)
		return db.ErrNoPost
	}

	if !isAuthor {
		pg.log.Dbg("user is not post author",
			"post_id", req.PostID, "user_id", hubUserID)
		return db.ErrNotPostAuthor
	}

	// Delete existing comments if requested
	if req.DeleteExistingComments {
		_, err = tx.Exec(ctx, `
			DELETE FROM post_comments WHERE post_id = $1
		`, req.PostID)
		if err != nil {
			pg.log.Err("failed to delete existing comments",
				"error", err, "post_id", req.PostID)
			return db.ErrInternal
		}
	}

	// Disable comments for the post
	_, err = tx.Exec(ctx, `
		UPDATE posts SET comments_enabled = false WHERE id = $1
	`, req.PostID)
	if err != nil {
		pg.log.Err("failed to disable comments",
			"error", err, "post_id", req.PostID)
		return db.ErrInternal
	}

	err = tx.Commit(context.Background())
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return db.ErrInternal
	}

	return nil
}

func (pg *PG) EnablePostComments(
	ctx context.Context,
	req hub.EnablePostCommentsRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Check if post exists and user is the author, then enable comments
	result, err := pg.pool.Exec(ctx, `
		UPDATE posts 
		SET comments_enabled = true 
		WHERE id = $1 AND author_id = $2
	`, req.PostID, hubUserID)
	if err != nil {
		pg.log.Err("failed to enable comments",
			"error", err, "post_id", req.PostID)
		return db.ErrInternal
	}

	if result.RowsAffected() == 0 {
		// Check if post exists at all
		var postExists bool
		err = pg.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)
		`, req.PostID).Scan(&postExists)
		if err != nil {
			pg.log.Err("failed to check post existence",
				"error", err, "post_id", req.PostID)
			return db.ErrInternal
		}

		if !postExists {
			pg.log.Dbg("post not found", "post_id", req.PostID)
			return db.ErrNoPost
		}

		pg.log.Dbg("user is not post author",
			"post_id", req.PostID, "user_id", hubUserID)
		return db.ErrNotPostAuthor
	}

	return nil
}

func (pg *PG) DeletePostComment(
	ctx context.Context,
	req hub.DeletePostCommentRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Delete comment only if the user is the post author
	result, err := pg.pool.Exec(ctx, `
		DELETE FROM post_comments 
		WHERE id = $1 AND post_id = $2 
		AND EXISTS(SELECT 1 FROM posts WHERE id = $2 AND author_id = $3)
	`, req.CommentID, req.PostID, hubUserID)
	if err != nil {
		pg.log.Err("failed to delete comment",
			"error", err, "comment_id", req.CommentID, "post_id", req.PostID)
		return db.ErrInternal
	}

	if result.RowsAffected() == 0 {
		// Check if post exists
		var postExists bool
		err = pg.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)
		`, req.PostID).Scan(&postExists)
		if err != nil {
			pg.log.Err("failed to check post existence",
				"error", err, "post_id", req.PostID)
			return db.ErrInternal
		}

		if !postExists {
			pg.log.Dbg("post not found", "post_id", req.PostID)
			return db.ErrNoPost
		}

		// Check if user is post author
		var isAuthor bool
		err = pg.pool.QueryRow(ctx, `
			SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1 AND author_id = $2)
		`, req.PostID, hubUserID).Scan(&isAuthor)
		if err != nil {
			pg.log.Err("failed to check post authorship",
				"error", err, "post_id", req.PostID)
			return db.ErrInternal
		}

		if !isAuthor {
			pg.log.Dbg("user is not post author",
				"post_id", req.PostID, "user_id", hubUserID)
			return db.ErrNotPostAuthor
		}

		// Comment not found, but that's okay per the API spec
		pg.log.Dbg("comment not found",
			"comment_id", req.CommentID, "post_id", req.PostID)
	}

	return nil
}

func (pg *PG) DeleteMyComment(
	ctx context.Context,
	req hub.DeleteMyCommentRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	// Delete comment only if the user is the comment author
	result, err := pg.pool.Exec(ctx, `
		DELETE FROM post_comments 
		WHERE id = $1 AND post_id = $2 AND author_id = $3
	`, req.CommentID, req.PostID, hubUserID)
	if err != nil {
		pg.log.Err("failed to delete my comment",
			"error", err, "comment_id", req.CommentID, "post_id", req.PostID)
		return db.ErrInternal
	}

	if result.RowsAffected() == 0 {
		// Comment not found, but that's okay per the API spec
		pg.log.Dbg(
			"my comment not found",
			"comment_id",
			req.CommentID,
			"post_id",
			req.PostID,
			"user_id",
			hubUserID,
		)
	}

	return nil
}
