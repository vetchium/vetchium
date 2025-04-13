package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) FollowUser(ctx context.Context, handle string) error {
	pg.log.Inf("Entered PG FollowUser", "handle", handle)

	// Get logged-in user ID
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get logged in hub user ID", "error", err)
		return err
	}

	// Get target user ID from handle
	var targetUserID string
	err = pg.pool.QueryRow(
		ctx,
		"SELECT id FROM hub_users WHERE handle = $1",
		handle,
	).Scan(&targetUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("target user not found", "handle", handle)
			return db.ErrNoHubUser
		}

		pg.log.Err("failed to get target user ID", "error", err)
		return err
	}

	// Insert into following_relationships (or ignore if already exists)
	_, err = pg.pool.Exec(ctx,
		`INSERT INTO following_relationships
		(consuming_hub_user_id, producing_hub_user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		loggedInUserID, targetUserID)
	if err != nil {
		pg.log.Err("failed to insert following relationship", "error", err)
		return err
	}

	return nil
}

func (pg *PG) UnfollowUser(ctx context.Context, handle string) error {
	pg.log.Inf("Entered PG UnfollowUser", "handle", handle)

	// Get logged-in user ID
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get logged in hub user ID", "error", err)
		return err
	}

	// Get target user ID from handle
	var targetUserID string
	err = pg.pool.QueryRow(
		ctx,
		"SELECT id FROM hub_users WHERE handle = $1",
		handle,
	).Scan(&targetUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("target user not found", "handle", handle)
			return db.ErrNoHubUser
		}

		pg.log.Err("failed to get target user ID", "error", err)
		return err
	}

	// Delete the following relationship if it exists
	_, err = pg.pool.Exec(ctx,
		`DELETE FROM following_relationships
		WHERE consuming_hub_user_id = $1 AND producing_hub_user_id = $2`,
		loggedInUserID, targetUserID)
	if err != nil {
		pg.log.Err("failed to delete following relationship", "error", err)
		return err
	}

	return nil
}

func (pg *PG) GetFollowStatus(
	ctx context.Context,
	handle string,
) (hub.FollowStatus, error) {
	pg.log.Inf("Entered PG GetFollowStatus", "handle", handle)

	// Get logged-in user ID
	loggedInUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get logged in hub user ID", "error", err)
		return hub.FollowStatus{}, err
	}

	// Get target user ID from handle
	var targetUserID string
	var userState string
	err = pg.pool.QueryRow(
		ctx,
		"SELECT id, state FROM hub_users WHERE handle = $1",
		handle,
	).Scan(&targetUserID, &userState)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			pg.log.Dbg("target user not found", "handle", handle)
			return hub.FollowStatus{}, db.ErrNoHubUser
		}

		pg.log.Err("failed to get target user ID", "error", err)
		return hub.FollowStatus{}, err
	}

	// Check if the logged-in user is following the target user
	var isFollowing bool
	err = pg.pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM following_relationships
			WHERE consuming_hub_user_id = $1 AND producing_hub_user_id = $2
		)`,
		loggedInUserID, targetUserID).Scan(&isFollowing)
	if err != nil {
		pg.log.Err("failed to check following status", "error", err)
		return hub.FollowStatus{}, err
	}

	// For now, we're not implementing blocking functionality
	// So isBlocked is always false
	isBlocked := false

	// Check if user can be followed - only active users can be followed
	canFollow := userState == string(hub.ActiveHubUserState)

	return hub.FollowStatus{
		IsFollowing: isFollowing,
		IsBlocked:   isBlocked,
		CanFollow:   canFollow,
	}, nil
}
