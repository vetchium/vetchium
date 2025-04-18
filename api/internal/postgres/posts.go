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

	tagIDs := make([]string, 0, len(addPostReq.NewTags))
	newTagsInsertQuery := `
WITH inserted AS (
    -- Attempt to insert the tag.
    -- If it succeeds, return the new id and the name.
    -- If it conflicts (name exists), DO NOTHING and return nothing from this CTE.
    INSERT INTO tags (name)
    VALUES ($1)  -- $1 is your tag name parameter
    ON CONFLICT (name) DO NOTHING
    RETURNING id, name
)
-- First, try to select the id from the 'inserted' CTE (if the insert succeeded).
SELECT id
FROM inserted
UNION ALL
-- If the 'inserted' CTE is empty (meaning ON CONFLICT DO NOTHING happened),
-- select the id from the main 'tags' table where the name matches.
-- The 'WHERE NOT EXISTS (SELECT 1 FROM inserted)' clause ensures this part
-- only runs if the INSERT was skipped.
SELECT t.id
FROM tags t
WHERE t.name = $1 AND NOT EXISTS (SELECT 1 FROM inserted)
LIMIT 1; -- Ensures only one row is returned in any case
`

	for _, tag := range addPostReq.NewTags {
		var newTagID string
		err = tx.QueryRow(
			addPostReq.Context,
			newTagsInsertQuery,
			tag,
		).Scan(&newTagID)
		if err != nil {
			pg.log.Err("failed to insert new tags", "error", err)
			return err
		}
		tagIDs = append(tagIDs, newTagID)
	}

	for _, tagID := range addPostReq.TagIDs {
		tagIDs = append(tagIDs, string(tagID))
	}

	tagsInsertQuery := `
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`
	for _, tagID := range tagIDs {
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

	err = tx.Commit(context.Background())
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (pg *PG) GetPost(req db.GetPostRequest) (common.Post, error) {
	query := `
		SELECT
			p.id,
			p.content,
			p.created_at,
			p.updated_at,
			hu.handle AS author_handle,
			hu.full_name AS author_full_name,
			COALESCE(
				(
					SELECT json_agg(t.name ORDER BY t.name)
					FROM post_tags pt
					JOIN tags t ON pt.tag_id = t.id
					WHERE pt.post_id = p.id
				),
				'[]'::json
			) AS tags_json
		FROM
			posts p
		JOIN
			hub_users hu ON p.author_id = hu.id
		WHERE
			p.id = $1
	`

	var post common.Post
	var authorHandle string
	var authorFullName string
	var tagsJSON []byte

	err := pg.pool.QueryRow(req.Context, query, req.PostID).Scan(
		&post.ID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&authorHandle,
		&authorFullName,
		&tagsJSON,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("post not found", "post_id", req.PostID)
			return common.Post{}, db.ErrNoPost
		}
		pg.log.Err("failed to query post", "error", err, "post_id", req.PostID)
		return common.Post{}, db.ErrInternal
	}

	var tags []string
	if err := json.Unmarshal(tagsJSON, &tags); err != nil {
		pg.log.Err("unmarshalling", "error", err, "json", string(tagsJSON))
		return common.Post{}, db.ErrInternal
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

	// Query to fetch posts for the determined targetUserID
	query := `
		SELECT
			p.id,
			p.content,
			p.author_id,
			p.created_at,
			p.updated_at,
			hu.handle AS author_handle,
			hu.full_name AS author_full_name,
			COALESCE(
				(
					SELECT json_agg(t.name)
					FROM post_tags pt
					JOIN tags t ON pt.tag_id = t.id
					WHERE pt.post_id = p.id
				),
				'[]'::json
			) AS tags_json
		FROM
			posts p
		JOIN
			hub_users hu ON p.author_id = hu.id
		WHERE
			p.author_id = $1 -- Filter by the target user ID
	`
	// Start with the target user ID
	args := []interface{}{targetUserID}
	argCounter := 2 // Start arg counter for pagination from $2

	if getUserPostsReq.PaginationKey != nil &&
		*getUserPostsReq.PaginationKey != "" {
		var paginationUpdatedAt time.Time
		var paginationID string
		err := pg.pool.QueryRow(
			ctx,
			"SELECT updated_at, id FROM posts WHERE id = $1",
			*getUserPostsReq.PaginationKey,
		).Scan(&paginationUpdatedAt, &paginationID)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Proceed without pagination clause (effectively first page)
			} else {
				pg.log.Err("failed to get pagination key", "error", err)
				return hub.GetUserPostsResponse{}, err
			}
		} else {
			query += ` AND (p.updated_at, p.id) < ($` + fmt.Sprintf("%d", argCounter) + `::timestamptz, $` + fmt.Sprintf("%d", argCounter+1) + `)`
			args = append(args, paginationUpdatedAt, paginationID)
			argCounter += 2
		}
	}

	query += ` ORDER BY p.updated_at DESC, p.id DESC LIMIT $` + fmt.Sprintf(
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

	posts := make([]common.Post, 0, getUserPostsReq.Limit)
	var lastPostID string

	for rows.Next() {
		var post common.Post
		var authorID string
		var authorHandle string
		var authorFullName string
		var tagsJSON []byte

		err := rows.Scan(
			&post.ID,
			&post.Content,
			&authorID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&authorHandle,
			&authorFullName,
			&tagsJSON,
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
