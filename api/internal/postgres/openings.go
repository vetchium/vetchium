package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type CreateOpeningReq struct {
	Title              string
	Positions          int
	JD                 string
	Recruiters         []string
	HiringManager      string
	CostCenterName     string
	EmployerNotes      *string
	LocationTitles     []string
	RemoteCountryCodes []vetchi.CountryCode
	RemoteTimezones    []vetchi.TimeZone
	OpeningType        string
	YoeMin             int
	YoeMax             int
	MinEducationLevel  *vetchi.EducationLevel
	Salary             *vetchi.Salary
	EmployerID         uuid.UUID
	CreatedBy          uuid.UUID
}

type GetOpeningReq struct {
	ID         string
	EmployerID uuid.UUID
}

type FilterOpeningsReq struct {
	States        []string
	PaginationKey *string
	Limit         int
	EmployerID    uuid.UUID
}

type UpdateOpeningReq struct {
	ID                 string
	Title              string
	Positions          int
	JD                 string
	Recruiters         []string
	HiringManager      string
	CostCenterName     string
	EmployerNotes      *string
	LocationTitles     []string
	RemoteCountryCodes []vetchi.CountryCode
	RemoteTimezones    []vetchi.TimeZone
	OpeningType        string
	YoeMin             int
	YoeMax             int
	MinEducationLevel  *vetchi.EducationLevel
	Salary             *vetchi.Salary
	EmployerID         uuid.UUID
	UpdatedBy          uuid.UUID
}

type GetOpeningWatchersReq struct {
	ID         string
	EmployerID uuid.UUID
}

type AddOpeningWatchersReq struct {
	ID         string
	Emails     []string
	EmployerID uuid.UUID
	AddedBy    uuid.UUID
}

type RemoveOpeningWatcherReq struct {
	ID         string
	Email      string
	EmployerID uuid.UUID
	RemovedBy  uuid.UUID
}

type ApproveOpeningStateChangeReq struct {
	ID         string
	EmployerID uuid.UUID
	ApprovedBy uuid.UUID
}

type RejectOpeningStateChangeReq struct {
	ID         string
	EmployerID uuid.UUID
	RejectedBy uuid.UUID
}

// CreateOpening creates a new opening
func (pg *PG) CreateOpening(
	ctx context.Context,
	req CreateOpeningReq,
) (uuid.UUID, error) {
	// TODO: Implement this
	return uuid.Nil, nil
}

// GetOpening gets an opening by ID
func (pg *PG) GetOpening(
	ctx context.Context,
	req GetOpeningReq,
) (vetchi.Opening, error) {
	// TODO: Implement this
	return vetchi.Opening{}, nil
}

// FilterOpenings filters openings based on the given criteria
func (pg *PG) FilterOpenings(
	ctx context.Context,
	req FilterOpeningsReq,
) ([]vetchi.Opening, error) {
	// TODO: Implement this
	return nil, nil
}

// UpdateOpening updates an existing opening
func (pg *PG) UpdateOpening(ctx context.Context, req UpdateOpeningReq) error {
	// TODO: Implement this
	return nil
}

// GetOpeningWatchers gets the watchers of an opening
func (pg *PG) GetOpeningWatchers(
	ctx context.Context,
	req GetOpeningWatchersReq,
) (vetchi.OpeningWatchers, error) {
	// TODO: Implement this
	return vetchi.OpeningWatchers{}, nil
}

// AddOpeningWatchers adds watchers to an opening
func (pg *PG) AddOpeningWatchers(
	ctx context.Context,
	req AddOpeningWatchersReq,
) error {
	// TODO: Implement this
	return nil
}

// RemoveOpeningWatcher removes a watcher from an opening
func (pg *PG) RemoveOpeningWatcher(
	ctx context.Context,
	req RemoveOpeningWatcherReq,
) error {
	// TODO: Implement this
	return nil
}

// ApproveOpeningStateChange approves a pending state change for an opening
func (pg *PG) ApproveOpeningStateChange(
	ctx context.Context,
	req ApproveOpeningStateChangeReq,
) error {
	// TODO: Implement this
	return nil
}

// RejectOpeningStateChange rejects a pending state change for an opening
func (pg *PG) RejectOpeningStateChange(
	ctx context.Context,
	req RejectOpeningStateChangeReq,
) error {
	// TODO: Implement this
	return nil
}

// Add error constants
var (
	ErrNoOpening            = errors.New("opening not found")
	ErrNoStateChangeWaiting = errors.New("no state change waiting")
)
