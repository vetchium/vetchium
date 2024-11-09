package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

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
