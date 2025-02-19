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
WITH verified_domains AS (
	SELECT DISTINCT d.domain_name
	FROM hub_users_official_emails hoe
	JOIN domains d ON d.id = hoe.domain_id
	JOIN hub_users hu ON hu.id = hoe.hub_user_id
	WHERE hu.handle = $1
	AND hoe.last_verified_at IS NOT NULL
)
SELECT hu.handle, hu.full_name, hu.short_bio, hu.long_bio,
	COALESCE(array_agg(vd.domain_name) FILTER (WHERE vd.domain_name IS NOT NULL), '{}') as verified_mail_domains
FROM hub_users hu
LEFT JOIN verified_domains vd ON true
WHERE hu.handle = $1
GROUP BY hu.handle, hu.full_name, hu.short_bio, hu.long_bio
`, handle).Scan(&bio.Handle, &bio.FullName, &bio.ShortBio, &bio.LongBio, &bio.VerifiedMailDomains)
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
