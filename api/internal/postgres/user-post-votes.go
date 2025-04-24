package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) UpvoteUserPost(
	ctx context.Context,
	upvoteReq hub.UpvoteUserPostRequest,
) error {
	userID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user id", "error", err)
		return db.ErrInternal
	}

	// Check if post exists and if user has already downvoted in a single query
	var exists bool
	var hasDownvoted bool
	err = pg.pool.QueryRow(ctx, `
		WITH post_check AS (
			SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1) AS exists
		),
		downvote_check AS (
			SELECT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = $1 AND user_id = $2 AND vote_value = -1
			) AS has_downvoted
		)
		SELECT post_check.exists, downvote_check.has_downvoted
		FROM post_check, downvote_check
	`, upvoteReq.PostID, userID).Scan(&exists, &hasDownvoted)
	if err != nil {
		pg.log.Err("failed", "error", err, "post_id", upvoteReq.PostID)
		return db.ErrInternal
	}

	if !exists {
		pg.log.Dbg("post does not exist", "post_id", upvoteReq.PostID)
		return db.ErrNonVoteableUserPost
	}

	if hasDownvoted {
		pg.log.Dbg("Already Downvoted",
			"post_id", upvoteReq.PostID,
			"user_id", userID)
		return db.ErrNonVoteableUserPost
	}

	// Insert or update vote
	_, err = pg.pool.Exec(ctx, `
		INSERT INTO post_votes (post_id, user_id, vote_value)
		VALUES ($1, $2, 1)
		ON CONFLICT (post_id, user_id) DO UPDATE
		SET vote_value = 1
	`, upvoteReq.PostID, userID)
	if err != nil {
		pg.log.Err("failed to register upvote",
			"error", err,
			"post_id", upvoteReq.PostID,
			"user_id", userID)
		return db.ErrInternal
	}

	pg.log.Dbg("upvoted", "post_id", upvoteReq.PostID, "user_id", userID)
	return nil
}

func (pg *PG) DownvoteUserPost(
	ctx context.Context,
	downvoteReq hub.DownvoteUserPostRequest,
) error {
	userID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user id", "error", err)
		return db.ErrInternal
	}

	// Check if post exists and if user has already upvoted in a single query
	var exists bool
	var hasUpvoted bool
	err = pg.pool.QueryRow(ctx, `
		WITH post_check AS (
			SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1) AS exists
		),
		upvote_check AS (
			SELECT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = $1 AND user_id = $2 AND vote_value = 1
			) AS has_upvoted
		)
		SELECT post_check.exists, upvote_check.has_upvoted
		FROM post_check, upvote_check
	`, downvoteReq.PostID, userID).Scan(&exists, &hasUpvoted)
	if err != nil {
		pg.log.Err("failed to downvote",
			"error", err,
			"post_id", downvoteReq.PostID,
			"user_id", userID)
		return db.ErrInternal
	}

	if !exists {
		pg.log.Dbg("post does not exist", "post_id", downvoteReq.PostID)
		return db.ErrNonVoteableUserPost
	}

	if hasUpvoted {
		pg.log.Dbg("already upvoted",
			"post_id", downvoteReq.PostID,
			"user_id", userID)
		return db.ErrNonVoteableUserPost
	}

	// Insert or update vote
	_, err = pg.pool.Exec(ctx, `
		INSERT INTO post_votes (post_id, user_id, vote_value)
		VALUES ($1, $2, -1)
		ON CONFLICT (post_id, user_id) DO UPDATE
		SET vote_value = -1
	`, downvoteReq.PostID, userID)
	if err != nil {
		pg.log.Err("failed to downvote",
			"error", err,
			"post_id", downvoteReq.PostID,
			"user_id", userID)
		return db.ErrInternal
	}

	pg.log.Dbg("downvoted", "post_id", downvoteReq.PostID, "user_id", userID)
	return nil
}

func (pg *PG) UnvoteUserPost(
	ctx context.Context,
	unvoteReq hub.UnvoteUserPostRequest,
) error {
	userID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user id", "error", err)
		return db.ErrInternal
	}

	// Check if post exists
	var exists bool
	err = pg.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM posts WHERE id = $1)
	`, unvoteReq.PostID).Scan(&exists)
	if err != nil {
		pg.log.Err("failed to check post exists",
			"error", err,
			"post_id", unvoteReq.PostID)
		return db.ErrInternal
	}

	if !exists {
		pg.log.Dbg("post does not exist", "post_id", unvoteReq.PostID)
		return db.ErrNonVoteableUserPost
	}

	// Delete any existing vote
	_, err = pg.pool.Exec(ctx, `
		DELETE FROM post_votes
		WHERE post_id = $1 AND user_id = $2
	`, unvoteReq.PostID, userID)
	if err != nil {
		pg.log.Err("failed to delete vote",
			"error", err,
			"post_id", unvoteReq.PostID,
			"user_id", userID)
		return db.ErrInternal
	}

	pg.log.Dbg("unvoted", "post_id", unvoteReq.PostID, "user_id", userID)
	return nil
}
