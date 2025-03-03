package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/employer"
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
),
connection_state AS (
	SELECT 
		CASE
			-- Check if they are connected
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE (requester_id = $2 AND requested_id = (SELECT id FROM target_user_id)
					OR requester_id = (SELECT id FROM target_user_id) AND requested_id = $2)
				AND state = 'COLLEAGUING_ACCEPTED'
			) THEN 'CONNECTED'
			
			-- Check if there's a pending request from logged-in user
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE requester_id = $2 
				AND requested_id = (SELECT id FROM target_user_id)
				AND state = 'COLLEAGUING_PENDING'
			) THEN 'REQUEST_SENT_PENDING'
			
			-- Check if there's a pending request to logged-in user
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE requester_id = (SELECT id FROM target_user_id)
				AND requested_id = $2
				AND state = 'COLLEAGUING_PENDING'
			) THEN 'REQUEST_RECEIVED_PENDING'
			
			-- Check if logged-in user rejected their request
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE requester_id = (SELECT id FROM target_user_id)
				AND requested_id = $2
				AND state = 'COLLEAGUING_REJECTED'
				AND rejected_by = $2
			) THEN 'REJECTED_BY_ME'
			
			-- Check if they rejected logged-in user's request
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE requester_id = $2
				AND requested_id = (SELECT id FROM target_user_id)
				AND state = 'COLLEAGUING_REJECTED'
				AND rejected_by = (SELECT id FROM target_user_id)
			) THEN 'REJECTED_BY_THEM'
			
			-- Check if logged-in user unlinked the connection
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE (requester_id = $2 AND requested_id = (SELECT id FROM target_user_id)
					OR requester_id = (SELECT id FROM target_user_id) AND requested_id = $2)
				AND state = 'COLLEAGUING_UNLINKED'
				AND unlinked_by = $2
			) THEN 'UNLINKED_BY_ME'
			
			-- Check if they unlinked the connection
			WHEN EXISTS (
				SELECT 1 FROM colleague_connections
				WHERE (requester_id = $2 AND requested_id = (SELECT id FROM target_user_id)
					OR requester_id = (SELECT id FROM target_user_id) AND requested_id = $2)
				AND state = 'COLLEAGUING_UNLINKED'
				AND unlinked_by = (SELECT id FROM target_user_id)
			) THEN 'UNLINKED_BY_THEM'
			
			-- If none of the above, check if they can be colleagues
			ELSE 
				CASE 
					WHEN is_colleaguable($2, (SELECT id FROM target_user_id), $3) THEN 'CAN_SEND_REQUEST'
					ELSE 'CANNOT_SEND_REQUEST'
				END
		END as state
)
SELECT hu.handle, hu.full_name, hu.short_bio, hu.long_bio,
	COALESCE(array_agg(vd.domain_name) FILTER (WHERE vd.domain_name IS NOT NULL), '{}') as verified_mail_domains,
	cs.state as colleague_connection_state
FROM hub_users hu
LEFT JOIN verified_domains vd ON true
CROSS JOIN connection_state cs
WHERE hu.handle = $1
GROUP BY hu.handle, hu.full_name, hu.short_bio, hu.long_bio, cs.state
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
		&bio.ColleagueConnectionState,
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

func (p *PG) GetEmployerViewBio(
	ctx context.Context,
	handle string,
) (employer.EmployerViewBio, error) {
	var bio employer.EmployerViewBio

	err := p.pool.QueryRow(
		ctx,
		`
WITH hub_user_info AS (
    SELECT 
        hu.handle,
        hu.full_name,
        hu.short_bio,
        hu.long_bio,
        COALESCE(array_agg(DISTINCT d.domain_name) FILTER (WHERE d.domain_name IS NOT NULL), '{}') as verified_mail_domains
    FROM hub_users hu
    LEFT JOIN hub_users_official_emails hoe ON hu.id = hoe.hub_user_id AND hoe.last_verified_at IS NOT NULL
    LEFT JOIN domains d ON hoe.domain_id = d.id
    WHERE hu.handle = $1
    GROUP BY hu.handle, hu.full_name, hu.short_bio, hu.long_bio
),
work_history_info AS (
    SELECT 
        jsonb_agg(
            jsonb_build_object(
                'id', wh.id,
                'employer_domain', d.domain_name,
                'employer_name', e.company_name,
                'title', wh.title,
                'start_date', wh.start_date,
                'end_date', wh.end_date,
                'description', wh.description
            ) ORDER BY wh.start_date DESC
        ) as work_history
    FROM hub_users hu
    JOIN work_history wh ON hu.id = wh.hub_user_id
    JOIN employers e ON e.id = wh.employer_id
    JOIN domains d ON d.employer_id = e.id AND d.id = (
        SELECT epd.domain_id 
        FROM employer_primary_domains epd 
        WHERE epd.employer_id = e.id
    )
    WHERE hu.handle = $1
)
SELECT 
    hui.*,
    COALESCE(whi.work_history, '[]'::jsonb) as work_history
FROM hub_user_info hui
LEFT JOIN work_history_info whi ON true
`,
		handle,
	).Scan(
		&bio.Handle,
		&bio.FullName,
		&bio.ShortBio,
		&bio.LongBio,
		&bio.VerifiedMailDomains,
		&bio.WorkHistory,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no hub user found", "handle", handle)
			return employer.EmployerViewBio{}, db.ErrNoHubUser
		}

		p.log.Err("failed to get bio", "error", err)
		return employer.EmployerViewBio{}, err
	}

	return bio, nil
}
