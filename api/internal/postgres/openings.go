package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) CreateOpening(
	ctx context.Context,
	createOpeningReq vetchi.CreateOpeningRequest,
) (string, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Error("failed to get orgUser from context")
		return "", db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var costCenterID uuid.UUID
	ccQuery := `SELECT id FROM cost_centers WHERE name = $1`
	err = tx.QueryRow(ctx, ccQuery, createOpeningReq.CostCenterName).
		Scan(&costCenterID)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Debug("CC not found", "name", createOpeningReq.CostCenterName)
			return "", db.ErrNoCostCenter
		}
		return "", err
	}

	todayOpeningsCountQuery := `SELECT COUNT(*) FROM openings WHERE employer_id = $1 AND created_at::date = CURRENT_DATE`
	var todayOpeningsCount int
	err = tx.QueryRow(ctx, todayOpeningsCountQuery, orgUser.EmployerID).
		Scan(&todayOpeningsCount)
	if err != nil {
		p.log.Error("failed to get today's openings count", "error", err)
		return "", err
	}

	// TODO: Check for max openings allowed per Employer and/or per Day.
	// Potentially this is where charging/pricing status codes will be returned
	var openingID string
	t := time.Now()
	openingID = fmt.Sprintf(
		"%d-%s-%d-%d",
		t.Year(),
		t.Format("Jan"),
		t.Day(),
		todayOpeningsCount+1,
	)
	p.log.Debug("generated opening ID", "id", openingID)

	query := `
INSERT INTO openings (id, title, positions, jd, hiring_manager, cost_center_id, employer_notes, remote_country_codes, remote_timezones, opening_type, yoe_min, yoe_max, min_education_level, salary_min, salary_max, salary_currency, current_state, approval_waiting_state, employer_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
RETURNING
    id
`
	err = tx.QueryRow(ctx, query, openingID, createOpeningReq.Title, createOpeningReq.Positions, createOpeningReq.JD, createOpeningReq.HiringManager, costCenterID, createOpeningReq.EmployerNotes, createOpeningReq.RemoteCountryCodes, createOpeningReq.RemoteTimezones, createOpeningReq.OpeningType, createOpeningReq.YoeMin, createOpeningReq.YoeMax, createOpeningReq.MinEducationLevel, createOpeningReq.Salary.MinAmount, createOpeningReq.Salary.MaxAmount, createOpeningReq.Salary.Currency, vetchi.DraftOpening, nil, orgUser.EmployerID).
		Scan(&openingID)
	if err != nil {
		p.log.Error("failed to create opening", "error", err)
		return "", err
	}

	if len(createOpeningReq.Recruiters) > 0 {
		// First verify all recruiters exist
		verifyRecruitersQuery := `
SELECT COUNT(*)
FROM (
    SELECT UNNEST($1::text[]) AS email
    EXCEPT
    SELECT email FROM org_users
    WHERE employer_id = $2
) AS invalid_recruiters`

		var invalidCount int
		err = tx.QueryRow(
			ctx,
			verifyRecruitersQuery,
			createOpeningReq.Recruiters,
			orgUser.EmployerID,
		).Scan(&invalidCount)
		if err != nil {
			p.log.Error("failed to verify recruiters", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Debug("invalid recruiters found", "count", invalidCount)
			return "", db.ErrInvalidRecruiter
		}

		// If all recruiters are valid, proceed with insertion
		insertRecruitersQuery := `
INSERT INTO opening_recruiters (opening_id, recruiter_id)
SELECT $1, id
FROM org_users
WHERE email = ANY($2)
`
		_, err = tx.Exec(
			ctx,
			insertRecruitersQuery,
			openingID,
			createOpeningReq.Recruiters,
		)
		if err != nil {
			p.log.Error("failed to insert recruiters", "error", err)
			return "", err
		}
	}

	if len(createOpeningReq.HiringTeam) > 0 {
		// TODO: Parse the vetchi handles and insert the hub_users(id) to the table
	}

	if len(createOpeningReq.LocationTitles) > 0 {
		// First verify all locations exist
		verifyLocationsQuery := `
SELECT COUNT(*)
FROM (
    SELECT UNNEST($1::text[]) AS title
    EXCEPT
    SELECT title FROM locations
) AS invalid_locations`

		var invalidCount int
		err = tx.QueryRow(
			ctx,
			verifyLocationsQuery,
			createOpeningReq.LocationTitles,
		).Scan(&invalidCount)
		if err != nil {
			p.log.Error("failed to verify locations", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Debug("invalid locations found", "count", invalidCount)
			return "", db.ErrInvalidLocation
		}

		// If all locations are valid, proceed with insertion
		locationQuery := `
INSERT INTO opening_locations (opening_id, location_id)
SELECT $1, l.id
FROM locations l
WHERE l.title = ANY($2)
`
		_, err = tx.Exec(
			ctx,
			locationQuery,
			openingID,
			createOpeningReq.LocationTitles,
		)
		if err != nil {
			p.log.Error("failed to insert locations", "error", err)
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return "", err
	}

	return openingID, nil
}

func (pg *PG) GetOpening(
	ctx context.Context,
	getOpeningReq vetchi.GetOpeningRequest,
) (vetchi.Opening, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		pg.log.Error("failed to get orgUser from context")
		return vetchi.Opening{}, db.ErrInternal
	}

	query := `
SELECT
    o.id,
    o.title,
    o.positions,
    o.jd,
    o.hiring_manager,
    cc.name as cost_center_name,
    o.employer_notes,
    o.remote_country_codes,
    o.remote_timezones,
    o.opening_type,
    o.yoe_min,
    o.yoe_max,
    o.min_education_level,
    o.salary_min,
    o.salary_max,
    o.salary_currency,
    o.current_state,
    o.approval_waiting_state,
    ARRAY_AGG(DISTINCT jsonb_build_object('email', r.email, 'name', r.name)) FILTER (WHERE r.email IS NOT NULL) as recruiters,
    ARRAY_AGG(DISTINCT l.title) FILTER (WHERE l.title IS NOT NULL) as locations,
    ARRAY_AGG(DISTINCT jsonb_build_object('handle', hu.handle, 'full_name', hu.full_name)) FILTER (WHERE hu.handle IS NOT NULL) as hiring_team
FROM
    openings o
    LEFT JOIN cost_centers cc ON o.cost_center_id = cc.id
    LEFT JOIN opening_recruiters or2 ON o.id = or2.opening_id
    LEFT JOIN org_users r ON or2.recruiter_id = r.id
    LEFT JOIN opening_locations ol ON o.id = ol.opening_id
    LEFT JOIN locations l ON ol.location_id = l.id
    LEFT JOIN opening_hiring_team oht ON o.id = oht.opening_id
    LEFT JOIN hub_users hu ON oht.hub_user_id = hu.id
WHERE
    o.id = $1
    AND o.employer_id = $2
GROUP BY
    o.id, o.title, o.positions, o.jd, o.hiring_manager, cc.name,
    o.employer_notes, o.remote_country_codes, o.remote_timezones,
    o.opening_type, o.yoe_min, o.yoe_max, o.min_education_level,
    o.salary_min, o.salary_max, o.salary_currency, o.current_state,
    o.approval_waiting_state`

	var opening vetchi.Opening
	var locations []string
	var recruitersJSON, hiringTeamJSON [][]byte

	err := pg.pool.QueryRow(ctx, query, getOpeningReq.ID, orgUser.EmployerID).
		Scan(
			&opening.ID,
			&opening.Title,
			&opening.Positions,
			&opening.JD,
			&opening.HiringManager,
			&opening.CostCenterName,
			&opening.EmployerNotes,
			&opening.RemoteCountryCodes,
			&opening.RemoteTimezones,
			&opening.OpeningType,
			&opening.YoeMin,
			&opening.YoeMax,
			&opening.MinEducationLevel,
			&opening.Salary.MinAmount,
			&opening.Salary.MaxAmount,
			&opening.Salary.Currency,
			&opening.CurrentState,
			&opening.ApprovalWaitingState,
			&recruitersJSON,
			&locations,
			&hiringTeamJSON,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return vetchi.Opening{}, db.ErrNoOpening
		}
		pg.log.Error("failed to scan opening", "error", err)
		return vetchi.Opening{}, err
	}

	// Parse recruiters JSON
	opening.Recruiters = make([]vetchi.OrgUserShort, 0, len(recruitersJSON))
	for _, recruiterBytes := range recruitersJSON {
		var recruiter struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := json.Unmarshal(recruiterBytes, &recruiter); err != nil {
			pg.log.Error("failed to unmarshal recruiter", "error", err)
			return vetchi.Opening{}, err
		}
		opening.Recruiters = append(opening.Recruiters, vetchi.OrgUserShort{
			Email: recruiter.Email,
			Name:  recruiter.Name,
		})
	}

	// TODO: Parse hiring team

	opening.LocationTitles = locations

	return opening, nil
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
