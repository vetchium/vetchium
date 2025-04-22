package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) GetMyHomeTimeline(
	ctx context.Context,
	req hub.GetMyHomeTimelineRequest,
) (hub.MyHomeTimeline, error) {
	pg.log.Dbg("Entered PG GetMyHomeTimeline")

	// Get the logged-in user from context
	hubUserIDStr, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("Failed to get hub user ID from context", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf(
			"failed to get hub user ID: %w",
			err,
		)
	}

	hubUserID, err := uuid.Parse(hubUserIDStr)
	if err != nil {
		pg.log.Err("Failed to parse hub user ID", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("invalid hub user ID: %w", err)
	}

	// Start a transaction
	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		pg.log.Err("Failed to begin transaction", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Check if the pagination key exists if provided
	if req.PaginationKey != nil && *req.PaginationKey != "" {
		var exists bool
		err = tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM posts WHERE id = $1
			)
		`, *req.PaginationKey).Scan(&exists)

		if err != nil {
			pg.log.Err("Failed to check pagination key", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}

		if !exists {
			pg.log.Err("Invalid pagination key", "key", *req.PaginationKey)
			return hub.MyHomeTimeline{}, db.ErrInvalidPaginationKey
		}
	}

	// Check if user already has a timeline
	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM hu_active_home_timelines WHERE hub_user_id = $1
		)
	`, hubUserID).Scan(&exists)

	if err != nil {
		pg.log.Err("Failed to check if timeline exists", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}

	// If user doesn't have a timeline yet, create one and refresh it
	if !exists {
		pg.log.Dbg("Creating new timeline for user", "hub_user_id", hubUserID)

		// Initialize the timeline entry with an old refresh timestamp
		// so the initial RefreshTimeline call fetches recent history.
		_, err = tx.Exec(ctx, `
			INSERT INTO hu_active_home_timelines
				(hub_user_id, last_refreshed_at, last_accessed_at)
			VALUES
				($1, NOW() - INTERVAL '101 days', NOW())
		`, hubUserID)

		if err != nil {
			pg.log.Err("Failed to create timeline entry", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}

		// Call RefreshTimeline function to populate the timeline
		_, err = tx.Exec(ctx, `SELECT RefreshTimeline($1)`, hubUserID)
		if err != nil {
			pg.log.Err("Failed to refresh new timeline", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}
	} else {
		// Just update the last_accessed_at timestamp
		_, err = tx.Exec(ctx, `
			UPDATE hu_active_home_timelines
			SET last_accessed_at = NOW()
			WHERE hub_user_id = $1
		`, hubUserID)

		if err != nil {
			pg.log.Err("Failed to update last_accessed_at", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}
	}

	// Get posts from the timeline using the view
	var query string
	var args []interface{}

	if req.PaginationKey != nil && *req.PaginationKey != "" {
		// Query with pagination
		query = `
			SELECT
				post_id, content, created_at, updated_at,
				author_handle, author_name, author_profile_pic_url, tags
			FROM hu_timeline_extended
			WHERE hub_user_id = $1 AND post_id < $2
			ORDER BY most_recent_activity DESC, post_id DESC
			LIMIT $3
		`
		args = []interface{}{hubUserID, *req.PaginationKey, req.Limit}
	} else {
		// Query without pagination
		query = `
			SELECT
				post_id, content, created_at, author_handle, author_name, author_profile_pic_url, tags
			FROM hu_timeline_extended
			WHERE hub_user_id = $1
			ORDER BY most_recent_activity DESC, post_id DESC
			LIMIT $2
		`
		args = []interface{}{hubUserID, req.Limit}
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("Failed to query timeline posts", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var posts []common.Post
	var paginationKey string

	for rows.Next() {
		var post common.Post
		var profilePicURL *string
		var tags []string

		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&post.AuthorHandle,
			&post.AuthorName,
			&profilePicURL,
			&tags,
		)
		if err != nil {
			pg.log.Err("Failed to scan post row", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}

		// Tags are already retrieved from the view
		post.Tags = tags

		// Update pagination key to the last post ID
		paginationKey = post.ID

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("Error while iterating posts", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
		pg.log.Err("Failed to commit transaction", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}

	// Only include paginationKey if we have the maximum number of posts
	if len(posts) < req.Limit {
		paginationKey = ""
	}

	return hub.MyHomeTimeline{
		Posts:         posts,
		PaginationKey: paginationKey,
	}, nil
}
