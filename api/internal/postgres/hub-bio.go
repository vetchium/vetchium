package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetBio(ctx context.Context, handle string) (hub.Bio, error) {
	var bio hub.Bio
	err := p.pool.QueryRow(ctx, `
SELECT handle, full_name, short_bio, long_bio
FROM hub_users
WHERE handle = $1
`, handle).Scan(&bio.Handle, &bio.FullName, &bio.ShortBio, &bio.LongBio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.log.Dbg("no hub user found", "handle", handle)
			return hub.Bio{}, db.ErrNoHubUser
		}

		p.log.Err("failed to get bio", "error", err)
		return hub.Bio{}, err
	}
	return bio, nil
}

func (p *PG) UpdateBio(ctx context.Context, bio hub.UpdateBioRequest) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return err
	}

	query := `
UPDATE hub_users
SET
    handle = COALESCE($1, handle),
    full_name = COALESCE($2, full_name),
    short_bio = COALESCE($3, short_bio),
    long_bio = COALESCE($4, long_bio)
WHERE id = $5
`
	_, err = p.pool.Exec(
		ctx,
		query,
		bio.Handle,
		bio.FullName,
		bio.ShortBio,
		bio.LongBio,
		hubUserID,
	)
	if err != nil {
		p.log.Err("failed to update bio", "error", err)
		return err
	}

	return nil
}
