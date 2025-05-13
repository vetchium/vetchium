package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
)

func (p *PG) AddEmployerPost(req db.AddEmployerPostRequest) error {
	orgUser, ok := req.Context.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	tx, err := p.pool.Begin(req.Context)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
INSERT INTO employer_posts (id, content, employer_id)
VALUES ($1, $2, $3)
`
	_, err = tx.Exec(
		req.Context,
		query,
		req.PostID,
		req.Content,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to insert employer post", "error", err)
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
			p.log.Err("failed to insert new tags", "error", err)
			return err
		}
		tagIDs = append(tagIDs, newTagID)
	}

	for _, tagID := range req.TagIDs {
		tagIDs = append(tagIDs, string(tagID))
	}

	tagsInsertQuery := `
INSERT INTO employer_post_tags (post_id, tag_id)
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
			p.log.Err("failed to insert to post_tags", "error", err)
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}

func (p *PG) UpdateEmployerPost(req db.UpdateEmployerPostRequest) error {
	return nil
}

func (p *PG) DeleteEmployerPost(ctx context.Context, postID string) error {
	return nil
}
