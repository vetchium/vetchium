package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) CheckHandleAvailability(
	ctx context.Context,
	handle common.Handle,
) (hub.CheckHandleAvailabilityResponse, error) {
	pg.log.Dbg("checking handle availability", "handle", handle)

	var isAvailableForSignup bool
	err := pg.pool.QueryRow(ctx, `
SELECT NOT EXISTS (
	SELECT 1 FROM hub_users
	WHERE LOWER(handle) = LOWER($1)
)
`, string(handle)).Scan(&isAvailableForSignup)
	if err != nil {
		pg.log.Err("failed to check handle availability", "error", err)
		return hub.CheckHandleAvailabilityResponse{}, db.ErrInternal
	}

	pg.log.Dbg("handle availability", "bool", isAvailableForSignup)
	return hub.CheckHandleAvailabilityResponse{
		IsAvailable: isAvailableForSignup,
	}, nil
}

func (pg *PG) SetHandle(
	ctx context.Context,
	handle common.Handle,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	_, err = pg.pool.Exec(ctx, `
UPDATE hub_users
SET handle = $1
WHERE id = $2
AND tier = $3
`, string(handle), hubUserID, string(hub.PaidHubUserTier))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// probably not a paid user
			pg.log.Dbg("user is not a paid hub user")
			return db.ErrUnpaidHubUser
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "unique_handle":
				pg.log.Err("duplicate handle detected", "error", err)
				return db.ErrDupHandle
			default:
				pg.log.Err("unexpected constraint violation", "error", err)
				return db.ErrInternal
			}
		}
		pg.log.Err("failed to set handle", "error", err)
		return db.ErrInternal
	}

	pg.log.Dbg("handle set", "handle", handle)
	return nil
}
