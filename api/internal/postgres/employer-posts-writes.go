package postgres

import (
	"context"
	"fmt"

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

	tagIDs := make([]string, 0, len(req.NewTags)+len(req.TagIDs))
	for _, tag := range req.NewTags {
		var newTagID string

		// Use a proper upsert pattern that handles concurrency safely:
		// 1. First try to insert (ignoring conflicts)
		// 2. Then select the ID (which will always exist after step 1)

		// Try to insert the tag, ignore if it already exists
		_, err = tx.Exec(
			req.Context,
			"INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO NOTHING",
			tag,
		)
		if err != nil {
			p.log.Err("failed to insert tag", "tag_name", tag, "error", err)
			return db.ErrInternal
		}

		// Now select the ID (guaranteed to exist after the INSERT above)
		err = tx.QueryRow(
			req.Context,
			"SELECT id FROM tags WHERE name = $1",
			tag,
		).Scan(&newTagID)
		if err != nil {
			p.log.Err(
				"failed to select tag ID after upsert",
				"tag_name",
				tag,
				"error",
				err,
			)
			return db.ErrInternal
		}

		if newTagID == "" {
			p.log.Err("resolved newTagID is empty string", "tag_name", tag)
			return db.ErrInternal
		}
		tagIDs = append(tagIDs, newTagID)
	}

	// Validate existing TagIDs from req.TagIDs using a single SQL query for existence check.
	if len(req.TagIDs) > 0 {
		// 1. Convert req.TagIDs to a slice of non-empty strings.
		providedTagIDStrings := make([]string, 0, len(req.TagIDs))
		for _, tagIDInstance := range req.TagIDs {
			sTagID := string(tagIDInstance)
			if sTagID != "" {
				providedTagIDStrings = append(providedTagIDStrings, sTagID)
			}
		}

		// Only proceed with DB check if there are actual tag ID strings to validate.
		if len(providedTagIDStrings) > 0 {
			var allProvidedTagsExist bool
			query := `
SELECT NOT EXISTS (
	SELECT 1
	FROM unnest($1::text[]) AS pid_text
	LEFT JOIN tags t ON t.id = pid_text::uuid
	WHERE t.id IS NULL
)
`
			err = tx.QueryRow(req.Context, query, providedTagIDStrings).
				Scan(&allProvidedTagsExist)
			if err != nil {
				p.log.Err("failed tag existence check", "error", err)
				return db.ErrInternal
			}

			if !allProvidedTagsExist {
				p.log.Dbg("one or more provided tag IDs do not exist")
				return db.ErrNoTag
			}
		}

		// All tags are valid, append them to the tagIDs list.
		for _, existingTagUUID := range req.TagIDs {
			tagIDs = append(tagIDs, string(existingTagUUID))
		}
	}

	tagsInsertQuery := `
INSERT INTO employer_post_tags (employer_post_id, tag_id)
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
			p.log.Err("failed to insert into employer_post_tags", "error", err)
			return db.ErrInternal
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
	return fmt.Errorf("not implemented yet")
}

func (p *PG) DeleteEmployerPost(ctx context.Context, postID string) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("failed to begin transaction", "error", err)
		return err
	}
	defer tx.Rollback(context.Background())

	res, err := tx.Exec(
		ctx,
		"DELETE FROM employer_posts WHERE id = $1 AND employer_id = $2",
		postID,
		orgUser.EmployerID,
	)
	if err != nil {
		p.log.Err("failed to delete employer post", "error", err)
		return err
	}

	if res.RowsAffected() == 0 {
		p.log.Dbg("employer post not found", "post_id", postID)
		return db.ErrNoEmployerPost
	}

	_, err = tx.Exec(
		ctx,
		"DELETE FROM employer_post_tags WHERE employer_post_id = $1",
		postID,
	)
	if err != nil {
		p.log.Err("failed to delete employer post tags", "error", err)
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
