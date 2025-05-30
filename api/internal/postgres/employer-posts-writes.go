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

	// Validate existing TagIDs if provided
	if len(req.TagIDs) > 0 {
		// Deduplicate tag IDs
		uniqueTagIDs := make(map[string]bool)
		var deduplicatedTagIDs []string
		for _, tagID := range req.TagIDs {
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
		err = tx.QueryRow(req.Context, validateTagsQuery, deduplicatedTagIDs).
			Scan(&validTagCount)
		if err != nil {
			p.log.Err("failed to validate tag IDs", "error", err)
			return err
		}

		if validTagCount != len(deduplicatedTagIDs) {
			p.log.Dbg(
				"invalid tag IDs provided",
				"expected",
				len(deduplicatedTagIDs),
				"found",
				validTagCount,
			)
			return db.ErrInvalidTagIDs
		}

		// Insert post-tag relationships for existing tags (using deduplicated IDs)
		tagsInsertQuery := `
INSERT INTO employer_post_tags (employer_post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING
`
		for _, tagID := range deduplicatedTagIDs {
			_, err = tx.Exec(
				req.Context,
				tagsInsertQuery,
				req.PostID,
				tagID,
			)
			if err != nil {
				p.log.Err(
					"failed to insert into employer_post_tags",
					"error",
					err,
				)
				return err
			}
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

	// TODO: Should we keep them for some time for compliance ?
	_, err = tx.Exec(
		ctx,
		"DELETE FROM employer_post_tags WHERE employer_post_id = $1",
		postID,
	)
	if err != nil {
		p.log.Err("failed to delete employer post tags", "error", err)
		return err
	}

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

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
