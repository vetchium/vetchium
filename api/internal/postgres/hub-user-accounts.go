package postgres

import (
	"context"

	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) SignupHubUser(ctx context.Context, email string) error {
	return nil
}

func (p *PG) ChangeEmailAddress(ctx context.Context, email string) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	_, err = p.pool.Exec(
		ctx,
		"UPDATE hub_users SET email = $1 WHERE id = $2",
		email,
		hubUserID,
	)
	if err != nil {
		p.log.Err("failed to update email", "error", err)
		return db.ErrInternal
	}

	p.log.Dbg("email updated successfully", "email", email)
	return nil
}
