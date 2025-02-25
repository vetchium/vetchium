package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (p *PG) FilterColleagues(
	ctx context.Context,
	req hub.FilterColleaguesRequest,
) ([]hub.HubUserShort, error) {
	p.log.Dbg("filtering colleagues", "request", req)

	// Get the caller's ID from the context
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		p.log.Err("failed to get hub user ID from context", "error", err)
		return nil, db.ErrInternal
	}

	colleagues := []hub.HubUserShort{}

	// We need to check both the directions of the relationships
	// (requester and requested)
	query := `
		SELECT hu.handle, hu.full_name as name, hu.short_bio
		FROM hub_users hu
		JOIN colleague_connections cc ON 
			(cc.requester_id = $1 AND cc.requested_id = hu.id) OR 
			(cc.requested_id = $1 AND cc.requester_id = hu.id)
		WHERE cc.state = $3
			AND (hu.full_name ILIKE $2 OR hu.handle ILIKE $2)
		LIMIT $4
	`

	rows, err := p.pool.Query(
		ctx,
		query,
		hubUserID,
		"%"+req.Prefix+"%",
		db.ColleagueAccepted,
		req.Limit,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.log.Dbg("no colleagues found", "request", req)
			return colleagues, nil
		}

		p.log.Err("failed to query colleagues", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colleague hub.HubUserShort
		err = rows.Scan(&colleague.Handle, &colleague.Name, &colleague.ShortBio)
		if err != nil {
			p.log.Err("failed to scan colleague", "error", err)
			return nil, err
		}
		colleagues = append(colleagues, colleague)
	}
	if err := rows.Err(); err != nil {
		p.log.Err("failed to iterate over colleagues", "error", err)
		return nil, err
	}

	p.log.Dbg("filtered colleagues", "count", len(colleagues))
	return colleagues, nil
}
