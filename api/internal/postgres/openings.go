package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// CreateOpening creates a new opening
func (pg *PG) CreateOpening(
	ctx context.Context,
	req db.CreateOpeningReq,
) (uuid.UUID, error) {
	// TODO: Implement this
	return uuid.Nil, nil
}

// GetOpening gets an opening by ID
func (pg *PG) GetOpening(
	ctx context.Context,
	req db.GetOpeningReq,
) (vetchi.Opening, error) {
	// TODO: Implement this
	return vetchi.Opening{}, nil
}

// FilterOpenings filters openings based on the given criteria
func (pg *PG) FilterOpenings(
	ctx context.Context,
	req db.FilterOpeningsReq,
) ([]vetchi.Opening, error) {
	// TODO: Implement this
	return nil, nil
}

// UpdateOpening updates an existing opening
func (pg *PG) UpdateOpening(
	ctx context.Context,
	req db.UpdateOpeningReq,
) error {
	// TODO: Implement this
	return nil
}

// GetOpeningWatchers gets the watchers of an opening
func (pg *PG) GetOpeningWatchers(
	ctx context.Context,
	req db.GetOpeningWatchersReq,
) (vetchi.OpeningWatchers, error) {
	// TODO: Implement this
	return vetchi.OpeningWatchers{}, nil
}

// AddOpeningWatchers adds watchers to an opening
func (pg *PG) AddOpeningWatchers(
	ctx context.Context,
	req db.AddOpeningWatchersReq,
) error {
	// TODO: Implement this
	return nil
}

// RemoveOpeningWatcher removes a watcher from an opening
func (pg *PG) RemoveOpeningWatcher(
	ctx context.Context,
	req db.RemoveOpeningWatcherReq,
) error {
	// TODO: Implement this
	return nil
}

// ApproveOpeningStateChange approves a pending state change for an opening
func (pg *PG) ApproveOpeningStateChange(
	ctx context.Context,
	req db.ApproveOpeningStateChangeReq,
) error {
	// TODO: Implement this
	return nil
}

// RejectOpeningStateChange rejects a pending state change for an opening
func (pg *PG) RejectOpeningStateChange(
	ctx context.Context,
	req db.RejectOpeningStateChangeReq,
) error {
	// TODO: Implement this
	return nil
}
