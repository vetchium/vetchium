package postgres

import "github.com/vetchium/vetchium/api/internal/db"

func (pg *PG) AddPost(req db.AddPostRequest) error {
	hubUserID, err := getHubUserID(req.Context)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	insQuery := `
INSERT INTO posts (id, content, author_id)
VALUES ($1, $2, $3)
`

	_, err = pg.pool.Exec(
		req.Context,
		insQuery,
		req.PostID,
		req.Content,
		hubUserID,
	)
	if err != nil {
		pg.log.Err("failed to insert post", "error", err)
		return err
	}

	return nil
}
