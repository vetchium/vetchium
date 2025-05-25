package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) GetPostComments(
	ctx context.Context,
	req hub.GetPostCommentsRequest,
) ([]hub.PostComment, error) {
	// Check if post exists
	var postExists bool
	err := pg.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)
	`, req.PostID).Scan(&postExists)
	if err != nil {
		pg.log.Err("failed to check post existence",
			"error", err, "post_id", req.PostID)
		return nil, db.ErrInternal
	}

	if !postExists {
		pg.log.Dbg("post not found", "post_id", req.PostID)
		return nil, db.ErrNoPost
	}

	// Build query with pagination
	query := `
		SELECT 
			pc.id,
			pc.content,
			hu.full_name,
			hu.handle,
			pc.created_at
		FROM post_comments pc
		JOIN hub_users hu ON pc.author_id = hu.id
		WHERE pc.post_id = $1
	`
	args := []interface{}{req.PostID}
	argCount := 1

	// Add pagination condition if provided
	if req.PaginationKey != "" {
		// For pagination, we need to find comments created before the pagination key comment
		// or with the same created_at but lower ID
		query += ` AND (
			pc.created_at < (SELECT created_at FROM post_comments WHERE id = $2)
			OR (
				pc.created_at = (SELECT created_at FROM post_comments WHERE id = $2)
				AND pc.id < $2
			)
		)`
		argCount++
		args = append(args, req.PaginationKey)
	}

	// Order by newest first, then by ID descending for tie-breaking
	query += ` ORDER BY pc.created_at DESC, pc.id DESC`

	// Add limit
	argCount++
	query += ` LIMIT $` + string(rune('0'+argCount))
	args = append(args, req.Limit)

	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		pg.log.Err("failed to query comments",
			"error", err, "post_id", req.PostID)
		return nil, db.ErrInternal
	}
	defer rows.Close()

	var comments []hub.PostComment
	for rows.Next() {
		var comment hub.PostComment
		var authorName string
		var authorHandle string

		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&authorName,
			&authorHandle,
			&comment.CreatedAt,
		)
		if err != nil {
			pg.log.Err("failed to scan comment row", "error", err)
			return nil, db.ErrInternal
		}

		comment.AuthorName = authorName
		comment.AuthorHandle = common.Handle(authorHandle)
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		pg.log.Err("error iterating comment rows", "error", err)
		return nil, db.ErrInternal
	}

	return comments, nil
}
