package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// CreateOpening creates a new opening
func (pg *PG) CreateOpening(
	ctx context.Context,
	createOpeningReq vetchi.CreateOpeningRequest,
) (uuid.UUID, error) {
	// TODO: Implement this
	return uuid.Nil, nil
}

// GetOpening gets an opening by ID
func (pg *PG) GetOpening(
	ctx context.Context,
	getOpeningReq vetchi.GetOpeningRequest,
) (vetchi.Opening, error) {
	// TODO: Implement this
	return vetchi.Opening{}, nil
}

// FilterOpenings filters openings based on the given criteria
func (pg *PG) FilterOpenings(
	ctx context.Context,
	filterOpeningsReq vetchi.FilterOpeningsRequest,
) ([]vetchi.Opening, error) {
	// TODO: Implement this
	return nil, nil
}

// UpdateOpening updates an existing opening
func (pg *PG) UpdateOpening(
	ctx context.Context,
	updateOpeningReq vetchi.UpdateOpeningRequest,
) error {
	// TODO: Implement this
	return nil
}

// GetOpeningWatchers gets the watchers of an opening
func (pg *PG) GetOpeningWatchers(
	ctx context.Context,
	getOpeningWatchersReq vetchi.GetOpeningWatchersRequest,
) (vetchi.OpeningWatchers, error) {
	// TODO: Implement this
	return vetchi.OpeningWatchers{}, nil
}

// AddOpeningWatchers adds watchers to an opening
func (pg *PG) AddOpeningWatchers(
	ctx context.Context,
	addOpeningWatchersReq vetchi.AddOpeningWatchersRequest,
) error {
	// TODO: Implement this
	return nil
}

// RemoveOpeningWatcher removes a watcher from an opening
func (pg *PG) RemoveOpeningWatcher(
	ctx context.Context,
	removeOpeningWatcherReq vetchi.RemoveOpeningWatcherRequest,
) error {
	// TODO: Implement this
	return nil
}

// ApproveOpeningStateChange approves a pending state change for an opening
func (pg *PG) ApproveOpeningStateChange(
	ctx context.Context,
	approveOpeningStateChangeReq vetchi.ApproveOpeningStateChangeRequest,
) error {
	// TODO: Implement this
	return nil
}

// RejectOpeningStateChange rejects a pending state change for an opening
func (pg *PG) RejectOpeningStateChange(
	ctx context.Context,
	rejectOpeningStateChangeReq vetchi.RejectOpeningStateChangeRequest,
) error {
	// TODO: Implement this
	return nil
}
