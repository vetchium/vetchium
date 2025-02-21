package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) ConnectColleague(ctx context.Context, handle string) error {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return err
	}

	// Get the target user's ID from their handle
	var targetUserID string
	err = p.pool.QueryRow(ctx, `
		SELECT id FROM hub_users WHERE handle = $1
	`, handle).Scan(&targetUserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("no hub user found", "handle", handle)
			return db.ErrNoHubUser
		}
		p.log.Err("failed to get target user ID", "error", err)
		return err
	}

	// Check if the connection is allowed using is_colleaguable function
	var isColleaguable bool
	err = p.pool.QueryRow(ctx, `
		SELECT is_colleaguable($1, $2, $3)
	`, loggedInUserID, targetUserID, vetchi.VerificationValidityDuration).Scan(&isColleaguable)
	if err != nil {
		p.log.Err("failed to check if colleaguable", "error", err)
		return err
	}

	if !isColleaguable {
		p.log.Dbg(
			"users cannot be connected",
			"requester",
			loggedInUserID,
			"requested",
			targetUserID,
		)
		return db.ErrNotColleaguable
	}

	// Create the colleague connection request
	_, err = p.pool.Exec(ctx, `
		INSERT INTO colleague_connections (
			requester_id,
			requested_id,
			state
		) VALUES (
			$1,
			$2,
			$3
		)
	`, loggedInUserID, targetUserID, db.ColleaguePending)
	if err != nil {
		p.log.Err("failed to create colleague connection", "error", err)
		return err
	}

	p.log.Dbg(
		"colleague connection request created",
		"requester",
		loggedInUserID,
		"requested",
		targetUserID,
	)
	return nil
}
