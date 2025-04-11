package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
)

func (pg *PG) AddPost(req db.AddPostRequest) error {
	hubUserID, err := getHubUserID(req.Context)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	tx, err := pg.pool.Begin(req.Context)
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
		req.Context,
		postsInsertQuery,
		req.PostID,
		req.Content,
		hubUserID,
	)
	if err != nil {
		pg.log.Err("failed to insert post", "error", err)
		return err
	}

	tagIDs := make([]string, 0, len(req.NewTags))
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

	for _, tag := range req.NewTags {
		var newTagID string
		err = tx.QueryRow(
			req.Context,
			newTagsInsertQuery,
			tag,
		).Scan(&newTagID)
		if err != nil {
			pg.log.Err("failed to insert new tags", "error", err)
			return err
		}
		tagIDs = append(tagIDs, newTagID)
	}

	for _, tagID := range req.TagIDs {
		tagIDs = append(tagIDs, string(tagID))
	}

	tagsInsertQuery := `
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`
	for _, tagID := range tagIDs {
		_, err = tx.Exec(
			req.Context,
			tagsInsertQuery,
			req.PostID,
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
