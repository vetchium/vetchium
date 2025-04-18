package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
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

	var args []interface{}

	query := `
		SELECT
			cc.id,
			hu.handle,
			hu.full_name as name,
			hu.short_bio
		FROM colleague_connections cc
		JOIN hub_users hu ON cc.requester_id = hu.id
		WHERE cc.requested_id = $1::uuid
		AND cc.state = $2`

	args = append(args, loggedInUserID, db.ColleaguePending)

	if req.PaginationKey != nil {
		query += fmt.Sprintf(` AND cc.id > $%d`, len(args)+1)
		args = append(args, *req.PaginationKey)
	}

	query += fmt.Sprintf(` ORDER BY cc.id ASC LIMIT $%d`, len(args)+1)
	args = append(args, req.Limit)

	hubUsers := []hub.HubUserShort{}
	var lastID string
	rows, err := p.pool.Query(ctx, query, args...)
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

	var args []interface{}
	args = append(args, loggedInUserID, db.ColleaguePending)

	if req.PaginationKey != nil {
		query += fmt.Sprintf(` AND cc.id < $%d`, len(args)+1)
		args = append(args, *req.PaginationKey)
	}

	query += fmt.Sprintf(` ORDER BY cc.id DESC LIMIT $%d`, len(args)+1)
	args = append(args, req.Limit)

	hubUsers := []hub.HubUserShort{}
	var lastID string
	rows, err := p.pool.Query(ctx, query, args...)
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

func (p *PG) ApproveColleague(ctx context.Context, handle string) error {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return err
	}

	// Single query that handles all cases: user not found, no pending request, and successful update
	query := `
		WITH requester AS (
			SELECT id, true as user_exists
			FROM hub_users
			WHERE handle = $1
		),
		connection_update AS (
			UPDATE colleague_connections cc
			SET state = $2
			FROM requester r
			WHERE cc.requester_id = r.id
			AND cc.requested_id = $3
			AND cc.state = $4
			RETURNING cc.requester_id as updated_id
		)
		SELECT
			COALESCE(r.user_exists, false) as user_exists,
			r.id as requester_id,
			cu.updated_id as connection_updated
		FROM requester r
		FULL OUTER JOIN connection_update cu ON true;
	`

	var userExists bool
	var requesterID, connectionUpdated *string
	err = p.pool.QueryRow(
		ctx,
		query,
		handle,
		db.ColleagueAccepted,
		loggedInUserID,
		db.ColleaguePending,
	).Scan(&userExists, &requesterID, &connectionUpdated)
	if err != nil {
		p.log.Err("failed to execute approve colleague query", "error", err)
		return err
	}

	if !userExists {
		p.log.Dbg("no hub user found", "handle", handle)
		return db.ErrNoHubUser
	}

	if connectionUpdated == nil {
		p.log.Dbg("no pending colleague request found",
			"handle", handle,
			"requested", loggedInUserID)
		return db.ErrNoApplication
	}

	p.log.Dbg("colleague connection approved",
		"requester", *requesterID,
		"requested", loggedInUserID)
	return nil
}

func (p *PG) RejectColleague(ctx context.Context, handle string) error {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return err
	}

	// Single query that handles all cases: user not found, no pending request, and successful update
	query := `
		WITH requester AS (
			SELECT id, true as user_exists
			FROM hub_users
			WHERE handle = $1
		),
		connection_update AS (
			UPDATE colleague_connections cc
			SET state = $2,
				rejected_by = $3,
				rejected_at = timezone('UTC', now()),
				updated_at = timezone('UTC', now())
			FROM requester r
			WHERE cc.requester_id = r.id
			AND cc.requested_id = $3
			AND cc.state = $4
			RETURNING cc.requester_id as updated_id
		)
		SELECT
			COALESCE(r.user_exists, false) as user_exists,
			r.id as requester_id,
			cu.updated_id as connection_updated
		FROM requester r
		FULL OUTER JOIN connection_update cu ON true;
	`

	var userExists bool
	var requesterID, connectionUpdated *string
	err = p.pool.QueryRow(
		ctx,
		query,
		handle,
		db.ColleagueRejected,
		loggedInUserID,
		db.ColleaguePending,
	).Scan(&userExists, &requesterID, &connectionUpdated)
	if err != nil {
		p.log.Err("failed to execute reject colleague query", "error", err)
		return err
	}

	if !userExists {
		p.log.Dbg("no hub user found", "handle", handle)
		return db.ErrNoHubUser
	}

	if connectionUpdated == nil {
		p.log.Dbg("no pending colleague request found",
			"handle", handle,
			"requested", loggedInUserID)
		return db.ErrNoApplication
	}

	p.log.Dbg("colleague connection rejected",
		"requester", *requesterID,
		"requested", loggedInUserID)
	return nil
}

func (p *PG) UnlinkColleague(ctx context.Context, handle string) error {
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in user ID", "error", err)
		return err
	}

	// Single query that handles all cases: user not found, no accepted connection, and successful update
	query := `
		WITH target_user AS (
			SELECT id, true as user_exists
			FROM hub_users
			WHERE handle = $1
		),
		connection_update AS (
			UPDATE colleague_connections cc
			SET state = $2,
				unlinked_by = $3,
				unlinked_at = timezone('UTC', now()),
				updated_at = timezone('UTC', now())
			FROM target_user t
			WHERE (
				(cc.requester_id = t.id AND cc.requested_id = $3) OR
				(cc.requester_id = $3 AND cc.requested_id = t.id)
			)
			AND cc.state = $4
			RETURNING 
				CASE 
					WHEN cc.requester_id = $3 THEN cc.requested_id
					ELSE cc.requester_id
				END as other_user_id
		)
		SELECT
			COALESCE(t.user_exists, false) as user_exists,
			t.id as target_id,
			cu.other_user_id as connection_updated
		FROM target_user t
		FULL OUTER JOIN connection_update cu ON true;
	`

	var userExists bool
	var targetID, connectionUpdated *string
	err = p.pool.QueryRow(
		ctx,
		query,
		handle,
		db.ColleagueUnlinked,
		loggedInUserID,
		db.ColleagueAccepted,
	).Scan(&userExists, &targetID, &connectionUpdated)

	if err != nil {
		p.log.Err("failed to execute unlink colleague query", "error", err)
		return err
	}

	if !userExists {
		p.log.Dbg("no hub user found", "handle", handle)
		return db.ErrNoHubUser
	}

	if connectionUpdated == nil {
		p.log.Dbg("no accepted colleague connection found",
			"handle", handle,
			"user", loggedInUserID)
		return db.ErrNoConnection
	}

	p.log.Dbg("colleague connection unlinked",
		"user", loggedInUserID,
		"other_user", *connectionUpdated)
	return nil
}
