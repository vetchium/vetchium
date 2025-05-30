package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) AddPost(addPostReq db.AddPostRequest) error {
	hubUserID, err := getHubUserID(addPostReq.Context)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tx, err := pg.pool.Begin(addPostReq.Context)
	if err != nil {
		pg.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	// Validate tag IDs if provided
	if len(addPostReq.TagIDs) > 0 {
		// Deduplicate tag IDs
		uniqueTagIDs := make(map[string]bool)
		var deduplicatedTagIDs []string
		for _, tagID := range addPostReq.TagIDs {
			tagIDStr := string(tagID)
			if !uniqueTagIDs[tagIDStr] {
				uniqueTagIDs[tagIDStr] = true
				deduplicatedTagIDs = append(deduplicatedTagIDs, tagIDStr)
			}
		}

		validateTagsQuery := `
SELECT COUNT(*) FROM tags WHERE id = ANY($1)
`
		var validTagCount int
		err = tx.QueryRow(addPostReq.Context, validateTagsQuery, deduplicatedTagIDs).
			Scan(&validTagCount)
		if err != nil {
			pg.log.Err("failed to validate tag IDs", "error", err)
			return err
		}

		if validTagCount != len(deduplicatedTagIDs) {
			pg.log.Dbg(
				"invalid tag IDs provided",
				"expected",
				len(deduplicatedTagIDs),
				"found",
				validTagCount,
			)
			return db.ErrInvalidTagIDs
		}
	}

	// Insert the post
	postsInsertQuery := `
INSERT INTO posts (id, content, author_id)
VALUES ($1, $2, $3)
`
	_, err = tx.Exec(
		addPostReq.Context,
		postsInsertQuery,
		addPostReq.PostID,
		addPostReq.Content,
		hubUserID,
	)
	if err != nil {
		pg.log.Err("failed to insert post", "error", err)
		return err
	}

	// Insert post-tag relationships for existing tags (using deduplicated IDs)
	if len(addPostReq.TagIDs) > 0 {
		// Deduplicate tag IDs again for insertion (reuse the logic)
		uniqueTagIDs := make(map[string]bool)
		var deduplicatedTagIDs []string
		for _, tagID := range addPostReq.TagIDs {
			tagIDStr := string(tagID)
			if !uniqueTagIDs[tagIDStr] {
				uniqueTagIDs[tagIDStr] = true
				deduplicatedTagIDs = append(deduplicatedTagIDs, tagIDStr)
			}
		}

		tagsInsertQuery := `
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`
		for _, tagID := range deduplicatedTagIDs {
			_, err = tx.Exec(
				addPostReq.Context,
				tagsInsertQuery,
				addPostReq.PostID,
				tagID,
			)
			if err != nil {
				pg.log.Err("failed to insert to post_tags", "error", err)
				return err
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (pg *PG) AddFTPost(addFTPostReq db.AddFTPostRequest) error {
	hubUserID, err := getHubUserID(addFTPostReq.Context)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tx, err := pg.pool.Begin(addFTPostReq.Context)
	if err != nil {
		pg.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	// Validate tag IDs if provided
	if len(addFTPostReq.TagIDs) > 0 {
		// Deduplicate tag IDs
		uniqueTagIDs := make(map[string]bool)
		var deduplicatedTagIDs []string
		for _, tagID := range addFTPostReq.TagIDs {
			tagIDStr := string(tagID)
			if !uniqueTagIDs[tagIDStr] {
				uniqueTagIDs[tagIDStr] = true
				deduplicatedTagIDs = append(deduplicatedTagIDs, tagIDStr)
			}
		}

		validateTagsQuery := `
SELECT COUNT(*) FROM tags WHERE id = ANY($1)
`
		var validTagCount int
		err = tx.QueryRow(addFTPostReq.Context, validateTagsQuery, deduplicatedTagIDs).
			Scan(&validTagCount)
		if err != nil {
			pg.log.Err("failed to validate tag IDs", "error", err)
			return err
		}

		if validTagCount != len(deduplicatedTagIDs) {
			pg.log.Dbg(
				"invalid tag IDs provided",
				"expected",
				len(deduplicatedTagIDs),
				"found",
				validTagCount,
			)
			return db.ErrInvalidTagIDs
		}
	}

	// Insert the post
	postsInsertQuery := `
INSERT INTO posts (id, content, author_id)
VALUES ($1, $2, $3)
`
	_, err = tx.Exec(
		addFTPostReq.Context,
		postsInsertQuery,
		addFTPostReq.PostID,
		addFTPostReq.Content,
		hubUserID,
	)
	if err != nil {
		pg.log.Err("failed to insert post", "error", err)
		return err
	}

	// Insert post-tag relationships for existing tags (using deduplicated IDs)
	if len(addFTPostReq.TagIDs) > 0 {
		// Deduplicate tag IDs again for insertion (reuse the logic)
		uniqueTagIDs := make(map[string]bool)
		var deduplicatedTagIDs []string
		for _, tagID := range addFTPostReq.TagIDs {
			tagIDStr := string(tagID)
			if !uniqueTagIDs[tagIDStr] {
				uniqueTagIDs[tagIDStr] = true
				deduplicatedTagIDs = append(deduplicatedTagIDs, tagIDStr)
			}
		}

		tagsInsertQuery := `
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`
		for _, tagID := range deduplicatedTagIDs {
			_, err = tx.Exec(
				addFTPostReq.Context,
				tagsInsertQuery,
				addFTPostReq.PostID,
				tagID,
			)
			if err != nil {
				pg.log.Err("failed to insert to post_tags", "error", err)
				return err
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (pg *PG) GetPost(req db.GetPostRequest) (hub.Post, error) {
	// Get the logged-in user's ID for voting status
	loggedInHubUserID, err := getHubUserID(req.Context)
	if err != nil {
		pg.log.Err("failed to get logged in hub user ID", "error", err)
		return hub.Post{}, err
	}

	query := `
		SELECT
			p.id,
			p.content,
			p.created_at,
			hu.handle AS author_handle,
			hu.full_name AS author_full_name,
			COALESCE(
				(
					SELECT json_agg(t.display_name ORDER BY t.display_name)
					FROM post_tags pt
					JOIN tags t ON pt.tag_id = t.id
					WHERE pt.post_id = p.id
				),
				'[]'::json
			) AS tags_json,
			p.author_id = $1 AS am_i_author,
			NOT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
			) AND p.author_id != $1 AS can_upvote,
			NOT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
			) AND p.author_id != $1 AS can_downvote,
			EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
				AND vote_value = 1
			) AS me_upvoted,
			EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
				AND vote_value = -1
			) AS me_downvoted,
			p.upvotes_count,
			p.downvotes_count,
			p.score,
			p.comments_enabled AS can_comment,
			(SELECT COUNT(*) FROM post_comments WHERE post_id = p.id)::int AS comments_count
		FROM
			posts p
		JOIN
			hub_users hu ON p.author_id = hu.id
		WHERE
			p.id = $2
	`

	var post hub.Post
	var authorHandle string
	var authorFullName string
	var tagsJSON []byte

	err = pg.pool.QueryRow(req.Context, query, loggedInHubUserID, req.PostID).
		Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&authorHandle,
			&authorFullName,
			&tagsJSON,
			&post.AmIAuthor,
			&post.CanUpvote,
			&post.CanDownvote,
			&post.MeUpvoted,
			&post.MeDownvoted,
			&post.UpvotesCount,
			&post.DownvotesCount,
			&post.Score,
			&post.CanComment,
			&post.CommentsCount,
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("post not found", "post_id", req.PostID)
			return hub.Post{}, db.ErrNoPost
		}
		pg.log.Err("failed to query post", "error", err, "post_id", req.PostID)
		return hub.Post{}, db.ErrInternal
	}

	var tags []string
	if err := json.Unmarshal(tagsJSON, &tags); err != nil {
		pg.log.Err("unmarshalling", "error", err, "json", string(tagsJSON))
		return hub.Post{}, db.ErrInternal
	}
	post.Tags = tags

	post.AuthorHandle = common.Handle(authorHandle)
	post.AuthorName = authorFullName

	return post, nil
}

func (pg *PG) GetUserPosts(
	ctx context.Context,
	getUserPostsReq hub.GetUserPostsRequest,
) (hub.GetUserPostsResponse, error) {
	loggedInHubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get logged in hub user ID", "error", err)
		return hub.GetUserPostsResponse{}, err
	}

	var targetUserID string
	if getUserPostsReq.Handle != nil && *getUserPostsReq.Handle != "" {
		// Handle provided: First, validate the handle and get the user ID.
		handle := string(*getUserPostsReq.Handle)
		err := pg.pool.QueryRow(ctx, "SELECT id FROM hub_users WHERE handle = $1", handle).
			Scan(&targetUserID)
		if err != nil {
			if err == pgx.ErrNoRows {
				pg.log.Dbg("non-existent handle", "handle", handle)
				return hub.GetUserPostsResponse{}, db.ErrNoHubUser
			}

			pg.log.Err("db error", "handle", handle, "error", err)
			return hub.GetUserPostsResponse{}, err
		}
	} else {
		// No handle provided: Use the logged-in user's ID.
		targetUserID = loggedInHubUserID
	}

	// Query to fetch posts with voting status and tags
	query := `
		SELECT
			p.id,
			p.content,
			p.created_at as created_at,
			hu.handle AS author_handle,
			hu.full_name AS author_full_name,
			COALESCE(
				(
					SELECT json_agg(t.display_name ORDER BY t.display_name)
					FROM post_tags pt
					JOIN tags t ON pt.tag_id = t.id
					WHERE pt.post_id = p.id
				),
				'[]'::json
			) AS tags_json,
			p.author_id = $1 AS am_i_author,
			NOT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
			) AND p.author_id != $1 AS can_upvote,
			NOT EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
			) AND p.author_id != $1 AS can_downvote,
			EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
				AND vote_value = 1
			) AS me_upvoted,
			EXISTS (
				SELECT 1 FROM post_votes
				WHERE post_id = p.id
				AND user_id = $1
				AND vote_value = -1
			) AS me_downvoted,
			p.upvotes_count,
			p.downvotes_count,
			p.score,
			p.comments_enabled AS can_comment,
			(SELECT COUNT(*) FROM post_comments WHERE post_id = p.id)::int AS comments_count
		FROM
			posts p
		JOIN
			hub_users hu ON p.author_id = hu.id
	`

	// Start with the logged in user ID
	args := []interface{}{loggedInHubUserID}
	argCounter := 2 // Start arg counter from $2

	// Build the WHERE clause based on conditions
	query += ` WHERE`

	// If a handle is provided, filter by that user's posts
	if getUserPostsReq.Handle != nil && *getUserPostsReq.Handle != "" {
		query += ` hu.handle = $` + fmt.Sprintf("%d", argCounter)
		args = append(args, string(*getUserPostsReq.Handle))
		argCounter++
	} else {
		// No handle provided - get logged in user's posts
		query += ` p.author_id = $1`
	}

	if getUserPostsReq.PaginationKey != nil &&
		*getUserPostsReq.PaginationKey != "" {
		var paginationCreatedAt time.Time
		var paginationID string
		err := pg.pool.QueryRow(
			ctx,
			"SELECT created_at, id FROM posts WHERE id = $1 LIMIT 1",
			*getUserPostsReq.PaginationKey,
		).Scan(&paginationCreatedAt, &paginationID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Proceed without pagination clause (effectively first page)
			} else {
				pg.log.Err("failed to get pagination key", "error", err)
				return hub.GetUserPostsResponse{}, err
			}
		} else {
			query += ` AND (p.created_at, p.id) < ($` + fmt.Sprintf("%d", argCounter) + `::timestamptz, $` + fmt.Sprintf("%d", argCounter+1) + `)`
			args = append(args, paginationCreatedAt, paginationID)
			argCounter += 2
		}
	}

	query += ` ORDER BY p.created_at DESC, p.id DESC LIMIT $` + fmt.Sprintf(
		"%d",
		argCounter,
	)
	args = append(args, getUserPostsReq.Limit)

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed", "error", err, "query", query, "args", args)
		return hub.GetUserPostsResponse{}, err
	}
	defer rows.Close()

	posts := make([]hub.Post, 0, getUserPostsReq.Limit)
	var lastPostID string

	for rows.Next() {
		var post hub.Post
		var tagsJSON []byte
		var authorHandle string
		var authorFullName string

		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.CreatedAt,
			&authorHandle,
			&authorFullName,
			&tagsJSON,
			&post.AmIAuthor,
			&post.CanUpvote,
			&post.CanDownvote,
			&post.MeUpvoted,
			&post.MeDownvoted,
			&post.UpvotesCount,
			&post.DownvotesCount,
			&post.Score,
			&post.CanComment,
			&post.CommentsCount,
		)
		if err != nil {
			pg.log.Err("failed to scan post row", "error", err)
			return hub.GetUserPostsResponse{}, err
		}

		var tags []string
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			pg.log.Err("JSON DB error", "error", err, "json", string(tagsJSON))
			return hub.GetUserPostsResponse{}, err
		}
		post.Tags = tags

		post.AuthorHandle = common.Handle(authorHandle)
		post.AuthorName = authorFullName

		posts = append(posts, post)
		lastPostID = string(post.ID)
	}

	if rows.Err() != nil {
		pg.log.Err("error iterating post rows", "error", rows.Err())
		return hub.GetUserPostsResponse{}, rows.Err()
	}

	var nextPaginationKey string
	if len(posts) == getUserPostsReq.Limit {
		nextPaginationKey = lastPostID
	}

	return hub.GetUserPostsResponse{
		Posts:         posts,
		PaginationKey: nextPaginationKey,
	}, nil
}
