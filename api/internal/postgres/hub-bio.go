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
