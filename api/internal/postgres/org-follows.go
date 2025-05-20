package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
)

func (p *PG) FollowOrg(ctx context.Context, domain string) error {
	// Get logged-in user ID
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in hub user ID", "error", err)
		return err
	}

	// Get employer ID from domain
	var employerID string
	err = p.pool.QueryRow(
		ctx,
		"SELECT employer_id FROM domains WHERE domain_name = $1",
		domain,
	).Scan(&employerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("domain not found", "domain", domain)
			return db.ErrNoDomain
		}
		p.log.Err("failed to get employer ID", "error", err)
		return err
	}

	// Insert into org_following_relationships (or ignore if already exists)
	_, err = p.pool.Exec(ctx,
		`INSERT INTO org_following_relationships
		(hub_user_id, employer_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`,
		hubUserID, employerID)
	if err != nil {
		p.log.Err("failed to insert org following relationship", "error", err)
		return err
	}

	return nil
}

func (p *PG) UnfollowOrg(ctx context.Context, domain string) error {
	// Get logged-in user ID
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get logged in hub user ID", "error", err)
		return err
	}

	// Get employer ID from domain
	var employerID string
	err = p.pool.QueryRow(
		ctx,
		"SELECT employer_id FROM domains WHERE domain_name = $1",
		domain,
	).Scan(&employerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("domain not found", "domain", domain)
			return db.ErrNoDomain
		}
		p.log.Err("failed to get employer ID", "error", err)
		return err
	}

	// Delete the following relationship if it exists
	_, err = p.pool.Exec(ctx,
		`DELETE FROM org_following_relationships
		WHERE hub_user_id = $1 AND employer_id = $2`,
		hubUserID, employerID)
	if err != nil {
		p.log.Err("failed to delete org following relationship", "error", err)
		return err
	}

	return nil
}
