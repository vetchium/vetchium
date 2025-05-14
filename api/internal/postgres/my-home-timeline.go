package postgres

import (
	"context"
	"database/sql" // Needed for sql.NullString and other nullable types from view
	"fmt"
	"time" // Added for time parsing

	"github.com/google/uuid"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common" // Assuming Handle might be here for casting
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

	// For pagination, we need the updated_at of the pagination key item
	var lastItemUpdatedAtString sql.NullString // Use sql.NullString for potential NULL from view
	if req.PaginationKey != nil && *req.PaginationKey != "" {
		err = tx.QueryRow(ctx, `
			SELECT updated_at FROM hu_timeline_extended
			WHERE hub_user_id = $1 AND item_id = $2
		`, hubUserID, *req.PaginationKey).Scan(&lastItemUpdatedAtString)

		if err != nil {
			if err == sql.ErrNoRows {
				pg.log.Dbg("Invalid pagination key", "key", *req.PaginationKey)
				return hub.MyHomeTimeline{}, db.ErrInvalidPaginationKey
			}
			pg.log.Err(
				"Failed to check pagination key and get updated_at",
				"error",
				err,
			)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}
		if !lastItemUpdatedAtString.Valid {
			// This case should ideally not happen if item_id is valid and in the view
			pg.log.Err(
				"Pagination key item has NULL updated_at",
				"key",
				*req.PaginationKey,
			)
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

	// The view hu_timeline_extended is already ORDER BY updated_at DESC, item_id DESC
	baseQuery := `
SELECT
    item_id, item_type, content, created_at, updated_at,
    author_handle, author_name, author_profile_pic_url,
    tags, upvotes_count, downvotes_count, score,
    me_upvoted, me_downvoted, can_upvote, can_downvote, am_i_author,
    employer_name, employer_id_internal, employer_domain_name
FROM hu_timeline_extended
WHERE hub_user_id = $1
`

	if req.PaginationKey != nil && *req.PaginationKey != "" &&
		lastItemUpdatedAtString.Valid {
		query = baseQuery + `AND (updated_at < $2 OR (updated_at = $2 AND item_id < $3)) ORDER BY updated_at DESC, item_id DESC LIMIT $4`
		args = []interface{}{
			hubUserID,
			lastItemUpdatedAtString.String,
			*req.PaginationKey,
			req.Limit,
		}
	} else {
		query = baseQuery + `ORDER BY updated_at DESC, item_id DESC LIMIT $2`
		args = []interface{}{hubUserID, req.Limit}
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err(
			"Failed to query timeline posts",
			"error",
			err,
			"query",
			query,
			"args",
			args,
		)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var userPosts []hub.Post
	var employerPosts []hub.EmployerPost
	var paginationKey string

	for rows.Next() {
		var itemID, itemTypeStr, content string
		var createdAtStr, updatedAtStr string // Read as string
		var authorHandle, authorName sql.NullString
		var authorProfilePicURL sql.NullString
		var tags []string // View already aggregates tags into an array of strings
		var upvotesCount, downvotesCount, score sql.NullInt32
		var meUpvoted, meDownvoted, canUpvote, canDownvote, amIAuthor sql.NullBool
		var employerName, employerIDInternal, employerDomainName sql.NullString

		err := rows.Scan(
			&itemID, &itemTypeStr, &content, &createdAtStr, &updatedAtStr,
			&authorHandle, &authorName, &authorProfilePicURL,
			&tags, &upvotesCount, &downvotesCount, &score,
			&meUpvoted, &meDownvoted, &canUpvote, &canDownvote, &amIAuthor,
			&employerName, &employerIDInternal, &employerDomainName,
		)
		if err != nil {
			pg.log.Err("Failed to scan timeline item row", "error", err)
			return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
		}

		// Common time parsing logic (assuming RFC3339 or similar from DB)
		parseTime := func(timeStr string) time.Time {
			if t, pErr := time.Parse(time.RFC3339, timeStr); pErr == nil {
				return t
			}
			// Fallback or error handling if parse fails. For now, return zero time.
			// Log an error here in a real scenario if parsing is critical.
			pg.log.Inf(
				"Failed to parse time string from DB, returning zero time",
				"time_string",
				timeStr,
			)
			return time.Time{}
		}

		itemType := common.TimelineItemType(itemTypeStr)

		if itemType == common.TimelineItemUserPost {
			userPost := hub.Post{
				ID:         itemID,
				Content:    content,
				Tags:       tags,
				AuthorName: authorName.String,
				AuthorHandle: common.Handle(
					authorHandle.String,
				), // Cast to common.Handle
				CreatedAt: parseTime(
					createdAtStr,
				), // Parse to time.Time
				UpvotesCount:   upvotesCount.Int32,
				DownvotesCount: downvotesCount.Int32,
				Score:          score.Int32,
				MeUpvoted:      meUpvoted.Bool,
				MeDownvoted:    meDownvoted.Bool,
				CanUpvote:      canUpvote.Bool,
				CanDownvote:    canDownvote.Bool,
				AmIAuthor:      amIAuthor.Bool,
				// ProfilePicURL needs to be handled if it's part of Post struct
			}
			userPosts = append(userPosts, userPost)
		} else if itemType == common.TimelineItemEmployerPost {
			employerPost := hub.EmployerPost{
				ID:                 itemID,
				Content:            content,
				Tags:               tags,
				EmployerName:       employerName.String,
				EmployerDomainName: employerDomainName.String,
				CreatedAt:          parseTime(createdAtStr), // Parse to time.Time
				UpdatedAt:          parseTime(updatedAtStr), // Parse to time.Time
			}
			employerPosts = append(employerPosts, employerPost)
		}

		paginationKey = itemID // The ID of the last processed item
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("Error while iterating timeline items", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
		pg.log.Err("Failed to commit transaction", "error", err)
		return hub.MyHomeTimeline{}, fmt.Errorf("database error: %w", err)
	}

	// Only include paginationKey if we have fetched up to the limit
	// The total number of items fetched is len(userPosts) + len(employerPosts)
	if (len(userPosts) + len(employerPosts)) < req.Limit {
		paginationKey = ""
	}

	return hub.MyHomeTimeline{
		Posts:         userPosts,
		EmployerPosts: employerPosts,
		PaginationKey: paginationKey,
	}, nil
}
