package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) GetOpening(
	ctx context.Context,
	getOpeningReq vetchi.GetOpeningRequest,
) (vetchi.Opening, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return vetchi.Opening{}, db.ErrInternal
	}

	query := `
SELECT
    o.id,
    o.title,
    o.positions,
    o.jd,
    jsonb_build_object('email', hm.email, 'name', hm.name, 'vetchi_handle', hu_hm.handle) AS hiring_manager,
    cc.cost_center_name,
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
    o.created_at,
    o.last_updated_at,
    jsonb_build_object('email', r.email, 'name', r.name, 'vetchi_handle', hu_r.handle) AS recruiter,
    ARRAY_AGG(DISTINCT l.title) FILTER (WHERE l.title IS NOT NULL) AS locations,
    ARRAY_AGG(DISTINCT jsonb_build_object('email', ht.email, 'name', ht.name, 'vetchi_handle', hu_ht.handle)) FILTER (WHERE ht.email IS NOT NULL) AS hiring_team
FROM
    openings o
    LEFT JOIN org_cost_centers cc ON o.cost_center_id = cc.id
    LEFT JOIN org_users hm ON o.hiring_manager = hm.id
    LEFT JOIN hub_users_official_emails hue_hm ON hm.email = hue_hm.official_email
    LEFT JOIN hub_users hu_hm ON hue_hm.hub_user_id = hu_hm.id
    LEFT JOIN org_users r ON o.recruiter = r.id
    LEFT JOIN hub_users_official_emails hue_r ON r.email = hue_r.official_email
    LEFT JOIN hub_users hu_r ON hue_r.hub_user_id = hu_r.id
    LEFT JOIN opening_locations ol ON o.id = ol.opening_id
    LEFT JOIN locations l ON ol.location_id = l.id
    LEFT JOIN opening_hiring_team oht ON o.id = oht.opening_id
    LEFT JOIN org_users ht ON oht.hiring_team_mate_id = ht.id
    LEFT JOIN hub_users_official_emails hue_ht ON ht.email = hue_ht.official_email
    LEFT JOIN hub_users hu_ht ON hue_ht.hub_user_id = hu_ht.id
WHERE
    o.id = $1
    AND o.employer_id = $2
GROUP BY
    o.id,
    o.title,
    o.positions,
    o.jd,
    hm.email,
    hm.name,
    hu_hm.handle,
    cc.cost_center_name,
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
    o.created_at,
    o.last_updated_at,
    r.email,
    r.name,
    hu_r.handle
`

	var opening vetchi.Opening
	var locations []string
	var hiringTeam []vetchi.OrgUserShort
	var recruiter, hiringManager vetchi.OrgUserShort
	var salary vetchi.Salary

	err := p.pool.QueryRow(ctx, query, getOpeningReq.ID, orgUser.EmployerID).
		Scan(
			&opening.ID,
			&opening.Title,
			&opening.Positions,
			&opening.JD,
			&hiringManager,
			&opening.CostCenterName,
			&opening.EmployerNotes,
			&opening.RemoteCountryCodes,
			&opening.RemoteTimezones,
			&opening.OpeningType,
			&opening.YoeMin,
			&opening.YoeMax,
			&opening.MinEducationLevel,
			&salary.MinAmount,
			&salary.MaxAmount,
			&salary.Currency,
			&opening.CurrentState,
			&opening.ApprovalWaitingState,
			&opening.CreatedAt,
			&opening.LastUpdatedAt,
			&recruiter,
			&locations,
			&hiringTeam,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return vetchi.Opening{}, db.ErrNoOpening
		}
		p.log.Err("failed to scan opening", "error", err)
		return vetchi.Opening{}, err
	}

	opening.Salary = &salary
	opening.LocationTitles = locations
	opening.HiringManager = hiringManager
	opening.Recruiter = recruiter
	opening.HiringTeam = hiringTeam

	return opening, nil
}

// FilterOpenings filters openings based on the given criteria
func (p *PG) FilterOpenings(
	ctx context.Context,
	filterOpeningsReq vetchi.FilterOpeningsRequest,
) ([]vetchi.OpeningInfo, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []vetchi.OpeningInfo{}, db.ErrInternal
	}

	query := `
WITH parsed_id AS (
    SELECT
        o.*,
        -- Split the ID into components and cast appropriately
        CAST(SPLIT_PART(id, '-', 1) AS INTEGER) as year,
        CAST(TO_CHAR(TO_DATE(SPLIT_PART(id, '-', 2), 'Mon'), 'MM') AS INTEGER) as month,
        CAST(SPLIT_PART(id, '-', 3) AS INTEGER) as day,
        CAST(SPLIT_PART(id, '-', 4) AS INTEGER) as sequence
    FROM openings o
    WHERE employer_id = $1
        AND current_state = ANY($2::opening_states[])
        AND created_at >= $3
        AND created_at <= $4
)
SELECT
    o.id,
    o.title,
    o.positions,
    0 as filled_positions, -- TODO Calculate filled positions
    o.opening_type,
    o.current_state,
    o.approval_waiting_state,
    o.created_at,
    o.last_updated_at,
    cc.cost_center_name,
    -- Recruiter details
    jsonb_build_object(
        'name', r.name,
        'email', r.email,
        'vetchi_handle', hu_r.handle
    ) as recruiter,
    -- Hiring Manager details
    jsonb_build_object(
        'name', hm.name,
        'email', hm.email,
        'vetchi_handle', hu_hm.handle
    ) as hiring_manager
FROM
    parsed_id o
    LEFT JOIN org_cost_centers cc ON o.cost_center_id = cc.id
    -- Join for Recruiter
    LEFT JOIN org_users r ON o.recruiter = r.id
    LEFT JOIN hub_users_official_emails hue_r ON r.email = hue_r.official_email
    LEFT JOIN hub_users hu_r ON hue_r.hub_user_id = hu_r.id
    -- Join for Hiring Manager
    LEFT JOIN org_users hm ON o.hiring_manager = hm.id
    LEFT JOIN hub_users_official_emails hue_hm ON hm.email = hue_hm.official_email
    LEFT JOIN hub_users hu_hm ON hue_hm.hub_user_id = hu_hm.id
WHERE
    -- Pagination logic
    CASE
        WHEN $5::TEXT IS NOT NULL AND $5::TEXT != '' THEN
            -- Parse the pagination key the same way
            (year, month, day, sequence) > (
                CAST(SPLIT_PART($5, '-', 1) AS INTEGER),
                CAST(TO_CHAR(TO_DATE(SPLIT_PART($5, '-', 2), 'Mon'), 'MM') AS INTEGER),
                CAST(SPLIT_PART($5, '-', 3) AS INTEGER),
                CAST(SPLIT_PART($5, '-', 4) AS INTEGER)
            )
        ELSE TRUE
    END
ORDER BY
    year,
    month,
    day,
    sequence
LIMIT $6
`

	rows, err := p.pool.Query(
		ctx,
		query,
		orgUser.EmployerID,
		filterOpeningsReq.StatesAsStrings(),
		filterOpeningsReq.FromDate,
		filterOpeningsReq.ToDate,
		filterOpeningsReq.PaginationKey,
		filterOpeningsReq.Limit,
	)
	if err != nil {
		p.log.Err("failed to query openings", "error", err)
		return []vetchi.OpeningInfo{}, err
	}

	openingInfos, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[vetchi.OpeningInfo],
	)
	if err != nil {
		p.log.Err("failed to collect rows", "error", err)
		return []vetchi.OpeningInfo{}, err
	}

	return openingInfos, nil
}

func (pg *PG) UpdateOpening(
	ctx context.Context,
	updateOpeningReq vetchi.UpdateOpeningRequest,
) error {
	// TODO: Implement this
	return nil
}

// GetOpeningWatchers gets the watchers of an opening
func (p *PG) GetOpeningWatchers(
	ctx context.Context,
	getOpeningWatchersReq vetchi.GetOpeningWatchersRequest,
) ([]vetchi.OrgUserShort, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []vetchi.OrgUserShort{}, db.ErrInternal
	}

	query := `
SELECT
    jsonb_build_object('email', ou.email, 'name', ou.name, 'vetchi_handle', ou.handle)
FROM
    opening_watchers ow
    LEFT JOIN org_users ou ON ow.watcher_id = ou.id
WHERE
    ow.employer_id = $1 AND ow.opening_id = $2
`

	rows, err := p.pool.Query(
		ctx,
		query,
		orgUser.EmployerID,
		getOpeningWatchersReq.OpeningID,
	)
	if err != nil {
		p.log.Err("failed to query opening watchers", "error", err)
		return []vetchi.OrgUserShort{}, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[vetchi.OrgUserShort])
}

func (p *PG) AddOpeningWatchers(
	ctx context.Context,
	addOpeningWatchersReq vetchi.AddOpeningWatchersRequest,
) error {
	// Expectations:
	// Invalid opening ID → db.ErrNoOpening
	// Invalid org user emails → db.ErrNoOrgUser
	// All users already watching → success (nil)
	// Mix of new and existing watchers → success (nil)
	// All new watchers → success (nil)
	// Database errors → db.ErrInternal
	// If adding watchers would exceed 25 watchers → db.ErrTooManyWatchers

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
WITH opening_check AS (
    -- First verify the opening exists and belongs to this employer
    SELECT EXISTS (
        SELECT 1 FROM openings
        WHERE id = $2 AND employer_id = $1
    ) as opening_exists
),
input_emails AS (
    -- Deduplicate input emails
    SELECT DISTINCT unnest($3::text[]) as email
),
org_users_to_add AS (
    -- Get org_user_ids for the given emails within the same org
    SELECT id, email
    FROM org_users
    WHERE email IN (SELECT email FROM input_emails)
    AND employer_id = $1
),
existing_watchers AS (
    -- Find which users are already watching
    SELECT watcher_id
    FROM opening_watchers
    WHERE employer_id = $1
    AND opening_id = $2
    AND watcher_id IN (SELECT id FROM org_users_to_add)
),
validation AS (
    SELECT
        (SELECT opening_exists FROM opening_check) as opening_exists,
        COUNT(*) = (SELECT COUNT(*) FROM input_emails) as all_emails_valid,
        (SELECT COUNT(*) FROM existing_watchers) = (SELECT COUNT(*) FROM org_users_to_add) as all_already_watching,
        (SELECT COUNT(*) FROM opening_watchers WHERE employer_id = $1 AND opening_id = $2) as current_watcher_count
    FROM org_users_to_add
),
insertion AS (
    INSERT INTO opening_watchers (employer_id, opening_id, watcher_id)
    SELECT $1, $2, org_users_to_add.id
    FROM org_users_to_add, validation
    WHERE validation.opening_exists
    AND validation.all_emails_valid
    AND validation.current_watcher_count + (SELECT COUNT(*) FROM org_users_to_add) <= 25
    AND NOT EXISTS (
        SELECT 1
        FROM opening_watchers ow
        WHERE ow.employer_id = $1
        AND ow.opening_id = $2
        AND ow.watcher_id = org_users_to_add.id
    )
    ON CONFLICT DO NOTHING
)
SELECT opening_exists, all_emails_valid, all_already_watching, current_watcher_count FROM validation;
`

	var openingExists, allEmailsValid, allAlreadyWatching bool
	var currentWatcherCount int
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUser.EmployerID,
		addOpeningWatchersReq.OpeningID,
		addOpeningWatchersReq.Emails,
	).Scan(&openingExists, &allEmailsValid, &allAlreadyWatching, &currentWatcherCount)
	if err != nil {
		p.log.Err("failed to add opening watchers", "error", err)
		return db.ErrInternal
	}

	if !openingExists {
		return db.ErrNoOpening
	}

	if !allEmailsValid {
		return db.ErrNoOrgUser
	}

	if currentWatcherCount+len(addOpeningWatchersReq.Emails) > 25 {
		return db.ErrTooManyWatchers
	}

	// If all users were already watching, that's still a success case
	return nil
}

func (p *PG) RemoveOpeningWatcher(
	ctx context.Context,
	removeOpeningWatcherReq vetchi.RemoveOpeningWatcherRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
DELETE FROM opening_watchers
WHERE employer_id = $1
    AND opening_id = $2
    AND watcher_id = (
        SELECT
            id
        FROM
            org_users
        WHERE
            email = $3
            AND employer_id = $1)
`

	_, err := p.pool.Exec(
		ctx,
		query,
		orgUser.EmployerID,
		removeOpeningWatcherReq.OpeningID,
		removeOpeningWatcherReq.Email,
	)
	if err != nil {
		p.log.Err("failed to remove opening watcher", "error", err)
		return db.ErrInternal
	}

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
