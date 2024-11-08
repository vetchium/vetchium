package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) CreateOpening(
	ctx context.Context,
	createOpeningReq vetchi.CreateOpeningRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Error("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}
	defer tx.Rollback(ctx)

	var costCenterID uuid.UUID
	ccQuery := `SELECT id FROM cost_centers WHERE name = $1`
	err = tx.QueryRow(ctx, ccQuery, createOpeningReq.CostCenterName).
		Scan(&costCenterID)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Debug("CC not found", "name", createOpeningReq.CostCenterName)
			return uuid.UUID{}, db.ErrNoCostCenter
		}
		return uuid.UUID{}, err
	}

	query := `
INSERT INTO openings (title, positions, jd, hiring_manager, cost_center_id, employer_notes, remote_country_codes, remote_timezones, opening_type, yoe_min, yoe_max, min_education_level, salary_min, salary_max, salary_currency, current_state, approval_waiting_state, employer_id, created_at, last_updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
RETURNING
    id
`
	var openingID uuid.UUID
	err = tx.QueryRow(ctx, query, createOpeningReq.Title, createOpeningReq.Positions, createOpeningReq.JD, createOpeningReq.HiringManager, costCenterID, createOpeningReq.EmployerNotes, createOpeningReq.RemoteCountryCodes, createOpeningReq.RemoteTimezones, createOpeningReq.OpeningType, createOpeningReq.YoeMin, createOpeningReq.YoeMax, createOpeningReq.MinEducationLevel, createOpeningReq.Salary.MinAmount, createOpeningReq.Salary.MaxAmount, createOpeningReq.Salary.Currency, vetchi.DraftOpening, nil, orgUser.ID).
		Scan(&openingID)
	if err != nil {
		p.log.Error("failed to create opening", "error", err)
		return uuid.UUID{}, err
	}

	if len(createOpeningReq.Recruiters) > 0 {
		insertRecruitersQuery := `
INSERT INTO opening_recruiters (opening_id, recruiter_id)
    VALUES ($1, (
            SELECT
                id
            FROM
                org_users
            WHERE
                email = $2))
RETURNING
    id
`
		for _, recruiter := range createOpeningReq.Recruiters {
			var recruiterID uuid.UUID
			err = tx.QueryRow(ctx, insertRecruitersQuery, openingID, string(recruiter)).
				Scan(&recruiterID)
			if err != nil {
				if err == sql.ErrNoRows {
					p.log.Debug("recruiter not found", "email", recruiter)
					return uuid.UUID{}, db.ErrNoRecruiter
				}
				p.log.Error("failed to insert recruiters", "error", err)
				return uuid.UUID{}, err
			}
		}
	}

	if len(createOpeningReq.HiringTeamMembers) > 0 {
		// TODO: Parse the vetchi handles and insert the hub_users(id) to the table
	}

	if len(createOpeningReq.LocationTitles) > 0 {
		locationQuery := `
INSERT INTO opening_locations (opening_id, location_id)
    VALUES ($1, (
            SELECT
                id
            FROM
                locations
            WHERE
                title = $2))
RETURNING
    id
`
		for _, location := range createOpeningReq.LocationTitles {
			var locationID uuid.UUID
			err = tx.QueryRow(ctx, locationQuery, openingID, location).
				Scan(&locationID)
			if err != nil {
				if err == sql.ErrNoRows {
					p.log.Debug("location not found", "title", location)
					return uuid.UUID{}, db.ErrNoLocation
				}

				p.log.Error("failed to insert locations", "error", err)
				return uuid.UUID{}, err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return uuid.UUID{}, err
	}

	return openingID, nil
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
