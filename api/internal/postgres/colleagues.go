package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/hub"
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

func (p *PG) GetMyColleagueApprovals(
	ctx context.Context,
	req hub.MyColleagueApprovalsRequest,
) (hub.MyColleagueApprovals, error) {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return hub.MyColleagueApprovals{}, err
	}

	query := `
		SELECT 
			cc.id,
			hu.handle,
			hu.full_name as name,
			hu.short_bio
		FROM colleague_connections cc
		JOIN hub_users hu ON cc.requester_id = hu.id
		WHERE cc.requested_id = $1
		AND cc.state = $2`

	if req.PaginationKey != nil {
		query += ` AND cc.id < $3`
	}

	query += ` ORDER BY cc.id DESC LIMIT $4`

	hubUsers := []hub.HubUserShort{}
	var lastID string
	rows, err := p.pool.Query(
		ctx,
		query,
		loggedInUserID,
		db.ColleaguePending,
		req.PaginationKey,
		req.Limit,
	)
	if err != nil {
		p.log.Err("failed to get my colleague approvals", "error", err)
		return hub.MyColleagueApprovals{}, err
	}

	for rows.Next() {
		var hubUser hub.HubUserShort
		var id string
		err = rows.Scan(&id, &hubUser.Handle, &hubUser.Name, &hubUser.ShortBio)
		if err != nil {
			p.log.Err("failed to scan hub user", "error", err)
			return hub.MyColleagueApprovals{}, err
		}
		lastID = id
		hubUsers = append(hubUsers, hubUser)
	}
	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over rows", "error", err)
		return hub.MyColleagueApprovals{}, err
	}

	return hub.MyColleagueApprovals{
		Approvals:     hubUsers,
		PaginationKey: lastID,
	}, nil
}

func (p *PG) GetMyColleagueSeeks(
	ctx context.Context,
	req hub.MyColleagueSeeksRequest,
) (hub.MyColleagueSeeks, error) {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return hub.MyColleagueSeeks{}, err
	}

	query := `
		SELECT 
			cc.id,
			hu.handle,
			hu.full_name as name,
			hu.short_bio
		FROM colleague_connections cc
		JOIN hub_users hu ON cc.requested_id = hu.id
		WHERE cc.requester_id = $1
		AND cc.state = $2`

	if req.PaginationKey != nil {
		query += ` AND cc.id < $3`
	}

	query += ` ORDER BY cc.id DESC LIMIT $4`

	hubUsers := []hub.HubUserShort{}
	var lastID string
	rows, err := p.pool.Query(
		ctx,
		query,
		loggedInUserID,
		db.ColleaguePending,
		req.PaginationKey,
		req.Limit,
	)
	if err != nil {
		p.log.Err("failed to get my colleague seeks", "error", err)
		return hub.MyColleagueSeeks{}, err
	}

	for rows.Next() {
		var hubUser hub.HubUserShort
		var id string
		err = rows.Scan(&id, &hubUser.Handle, &hubUser.Name, &hubUser.ShortBio)
		if err != nil {
			p.log.Err("failed to scan hub user", "error", err)
			return hub.MyColleagueSeeks{}, err
		}
		lastID = id
		hubUsers = append(hubUsers, hubUser)
	}
	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over rows", "error", err)
		return hub.MyColleagueSeeks{}, err
	}

	return hub.MyColleagueSeeks{
		Seeks:         hubUsers,
		PaginationKey: lastID,
	}, nil
}
