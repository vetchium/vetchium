package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetTimelinesToRefresh calls the GetOldestUnrefreshedActiveTimelines PostgreSQL function
// to retrieve hub_user_ids of timelines that need refreshing
func (p *PG) GetTimelinesToRefresh(
	ctx context.Context,
	limit int,
) ([]uuid.UUID, error) {
	p.log.Dbg("getting timelines to refresh", "limit", limit)

	rows, err := p.pool.Query(ctx, `
		SELECT * FROM GetOldestUnrefreshedActiveTimelines($1)
	`, limit)
	if err != nil {
		p.log.Err("failed to get timelines to refresh", "error", err)
		return nil, fmt.Errorf("failed to get timelines to refresh: %w", err)
	}
	defer rows.Close()

	var timelineIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			p.log.Err("failed to scan timeline ID", "error", err)
			return nil, fmt.Errorf("failed to scan timeline ID: %w", err)
		}
		timelineIDs = append(timelineIDs, id)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("error iterating timeline IDs", "error", err)
		return nil, fmt.Errorf("error iterating timeline IDs: %w", err)
	}

	p.log.Dbg("found timelines to refresh", "count", len(timelineIDs))
	return timelineIDs, nil
}

// RefreshTimeline calls the RefreshTimeline PostgreSQL function to update
// a user's home timeline with recent posts from followed users and self
func (p *PG) RefreshTimeline(
	ctx context.Context,
	hubUserID uuid.UUID,
) error {
	p.log.Dbg("refreshing timeline", "hub_user_id", hubUserID)

	_, err := p.pool.Exec(ctx, `
		SELECT RefreshTimeline($1)
	`, hubUserID)
	if err != nil {
		p.log.Err(
			"failed to refresh timeline",
			"hub_user_id",
			hubUserID,
			"error",
			err,
		)
		return fmt.Errorf(
			"failed to refresh timeline for user %s: %w",
			hubUserID,
			err,
		)
	}

	p.log.Dbg("timeline refreshed successfully", "hub_user_id", hubUserID)
	return nil
}
