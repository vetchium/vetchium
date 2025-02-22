package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) GetBio(ctx context.Context, handle string) (hub.Bio, error) {
	var bio hub.Bio
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return hub.Bio{}, err
	}

	err = p.pool.QueryRow(
		ctx,
		`
WITH verified_domains AS (
	SELECT DISTINCT d.domain_name
	FROM hub_users_official_emails hoe
	JOIN domains d ON d.id = hoe.domain_id
	JOIN hub_users hu ON hu.id = hoe.hub_user_id
	WHERE hu.handle = $1
	AND hoe.last_verified_at IS NOT NULL
),
target_user_id AS (
	SELECT id FROM hub_users WHERE handle = $1
)
SELECT hu.handle, hu.full_name, hu.short_bio, hu.long_bio,
	COALESCE(array_agg(vd.domain_name) FILTER (WHERE vd.domain_name IS NOT NULL), '{}') as verified_mail_domains,
	is_colleaguable($2, (SELECT id FROM target_user_id), $3) as is_colleaguable,
	are_colleagues($2, (SELECT id FROM target_user_id)) as is_colleague
FROM hub_users hu
LEFT JOIN verified_domains vd ON true
WHERE hu.handle = $1
GROUP BY hu.handle, hu.full_name, hu.short_bio, hu.long_bio
`,
		handle,
		loggedInUserID,
		vetchi.VerificationValidityDuration,
	).Scan(
		&bio.Handle,
		&bio.FullName,
		&bio.ShortBio,
		&bio.LongBio,
		&bio.VerifiedMailDomains,
		&bio.IsColleaguable,
		&bio.IsColleague,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
		// Check if this is a unique constraint violation on handle
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" &&
				pgerr.ConstraintName == "hub_users_handle_unique" {
				p.log.Dbg("duplicate handle", "handle", bio.Handle)
				return db.ErrDupHandle
			}
		}
		p.log.Err("failed to update bio", "error", err)
		return err
	}

	return nil
}
