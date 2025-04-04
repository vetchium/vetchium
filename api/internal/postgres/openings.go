package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func (p *PG) GetOpening(
	ctx context.Context,
	getOpeningReq employer.GetOpeningRequest,
) (employer.Opening, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.Opening{}, db.ErrInternal
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
    o.state,
    o.created_at,
    o.last_updated_at,
    jsonb_build_object('email', r.email, 'name', r.name, 'vetchi_handle', hu_r.handle) AS recruiter,
    ARRAY_AGG(DISTINCT l.title) FILTER (WHERE l.title IS NOT NULL) AS locations,
    ARRAY_AGG(DISTINCT jsonb_build_object('email', ht.email, 'name', ht.name, 'vetchi_handle', hu_ht.handle)) FILTER (WHERE ht.email IS NOT NULL) AS hiring_team,
    ARRAY_AGG(DISTINCT jsonb_build_object('id', ot.id, 'name', ot.name)) FILTER (WHERE ot.id IS NOT NULL) AS tags
FROM
    openings o
    LEFT JOIN org_cost_centers cc ON o.cost_center_id = cc.id
    LEFT JOIN org_users hm ON o.hiring_manager = hm.id
    LEFT JOIN hub_users_official_emails hue_hm ON hm.email = hue_hm.official_email
    LEFT JOIN hub_users hu_hm ON hue_hm.hub_user_id = hu_hm.id
    LEFT JOIN org_users r ON o.recruiter = r.id
    LEFT JOIN hub_users_official_emails hue_r ON r.email = hue_r.official_email
    LEFT JOIN hub_users hu_r ON hue_r.hub_user_id = hu_r.id
    LEFT JOIN opening_locations ol ON o.id = ol.opening_id AND o.employer_id = ol.employer_id
    LEFT JOIN locations l ON ol.location_id = l.id
    LEFT JOIN opening_hiring_team oht ON o.id = oht.opening_id AND o.employer_id = oht.employer_id
    LEFT JOIN org_users ht ON oht.hiring_team_mate_id = ht.id
    LEFT JOIN hub_users_official_emails hue_ht ON ht.email = hue_ht.official_email
    LEFT JOIN hub_users hu_ht ON hue_ht.hub_user_id = hu_ht.id
    LEFT JOIN opening_tag_mappings otm ON o.id = otm.opening_id AND o.employer_id = otm.employer_id
    LEFT JOIN opening_tags ot ON otm.tag_id = ot.id
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
    o.state,
    o.created_at,
    o.last_updated_at,
    r.email,
    r.name,
    hu_r.handle
`

	var opening employer.Opening
	var locations []string
	var hiringTeam []employer.OrgUserShort
	var recruiter, hiringManager employer.OrgUserShort
	var salary common.Salary
	var minAmount, maxAmount *float64
	var currencyStr *string
	var tags []common.OpeningTag

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
			&minAmount,
			&maxAmount,
			&currencyStr,
			&opening.State,
			&opening.CreatedAt,
			&opening.LastUpdatedAt,
			&recruiter,
			&locations,
			&hiringTeam,
			&tags,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return employer.Opening{}, db.ErrNoOpening
		}
		p.log.Err("failed to scan opening", "error", err)
		return employer.Opening{}, err
	}

	if minAmount != nil {
		salary.MinAmount = *minAmount
	}
	if maxAmount != nil {
		salary.MaxAmount = *maxAmount
	}
	if currencyStr != nil {
		salary.Currency = common.Currency(*currencyStr)
	}

	if minAmount != nil && maxAmount != nil && currencyStr != nil {
		opening.Salary = &salary
	}

	opening.LocationTitles = locations
	opening.HiringTeam = hiringTeam
	opening.HiringManager = hiringManager
	opening.Recruiter = recruiter
	opening.Tags = tags

	return opening, nil
}

// FilterOpenings filters openings based on the given criteria
func (p *PG) FilterOpenings(
	ctx context.Context,
	filterOpeningsReq employer.FilterOpeningsRequest,
) ([]employer.OpeningInfo, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []employer.OpeningInfo{}, db.ErrInternal
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
        AND state = ANY($2::opening_states[])
        AND created_at >= $3
        AND created_at < $4
)
SELECT
    o.id,
    o.title,
    o.positions,
    0 as filled_positions, -- TODO Calculate filled positions
    o.opening_type,
    o.state,
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
		return []employer.OpeningInfo{}, err
	}

	openingInfos, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[employer.OpeningInfo],
	)
	if err != nil {
		p.log.Err("failed to collect rows", "error", err)
		return []employer.OpeningInfo{}, err
	}

	return openingInfos, nil
}

func (pg *PG) UpdateOpening(
	ctx context.Context,
	updateOpeningReq employer.UpdateOpeningRequest,
) error {
	// TODO: Implement this
	return nil
}

// GetOpeningWatchers gets the watchers of an opening
func (p *PG) GetOpeningWatchers(
	ctx context.Context,
	getOpeningWatchersReq employer.GetOpeningWatchersRequest,
) ([]employer.OrgUserShort, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return []employer.OrgUserShort{}, db.ErrInternal
	}

	query := `
WITH opening_check AS (
    SELECT EXISTS (
        SELECT 1 
        FROM openings 
        WHERE id = $2 AND employer_id = $1
    ) as exists
),
watchers AS (
    SELECT
        ou.email,
        ou.name,
        COALESCE(hu.handle, '') as vetchi_handle
    FROM
        opening_watchers ow
        LEFT JOIN org_users ou ON ow.watcher_id = ou.id
        LEFT JOIN hub_users_official_emails hue ON ou.email = hue.official_email
        LEFT JOIN hub_users hu ON hue.hub_user_id = hu.id
    WHERE
        ow.employer_id = $1 AND ow.opening_id = $2
)
SELECT
    oc.exists as opening_exists,
    COALESCE(
        (SELECT json_agg(w.* ORDER BY w.email)
         FROM watchers w),
        '[]'
    ) as watchers
FROM opening_check oc;
`

	var openingExists bool
	var watchersJSON []byte
	err := p.pool.QueryRow(
		ctx,
		query,
		orgUser.EmployerID,
		getOpeningWatchersReq.OpeningID,
	).Scan(&openingExists, &watchersJSON)
	if err != nil {
		p.log.Err("failed to query opening watchers", "error", err)
		return []employer.OrgUserShort{}, db.ErrInternal
	}

	if !openingExists {
		return []employer.OrgUserShort{}, db.ErrNoOpening
	}

	var watchers []employer.OrgUserShort
	if err := json.Unmarshal(watchersJSON, &watchers); err != nil {
		p.log.Err("failed to unmarshal watchers", "error", err)
		return []employer.OrgUserShort{}, db.ErrInternal
	}

	return watchers, nil
}

func (p *PG) AddOpeningWatchers(
	ctx context.Context,
	addOpeningWatchersReq employer.AddOpeningWatchersRequest,
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
	removeOpeningWatcherReq employer.RemoveOpeningWatcherRequest,
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

func (p *PG) ChangeOpeningState(
	ctx context.Context,
	changeOpeningStateReq employer.ChangeOpeningStateRequest,
) error {
	const (
		resultNoOpening     = "no_opening"
		resultStateMismatch = "state_mismatch"
		resultUpdated       = "updated"
	)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
WITH opening_check AS (
    SELECT state
    FROM openings
    WHERE id = $1 AND employer_id = $2
),
state_update AS (
    UPDATE openings
    SET state = $3
    WHERE id = $1
        AND employer_id = $2
        AND state = $4
    RETURNING true as updated
)
SELECT 
    CASE
        WHEN NOT EXISTS (SELECT 1 FROM opening_check) THEN $5
        WHEN NOT EXISTS (SELECT 1 FROM state_update) THEN $6
        ELSE $7
    END as result;
`

	var result string
	err := p.pool.QueryRow(
		ctx,
		query,
		changeOpeningStateReq.OpeningID,
		orgUser.EmployerID,
		changeOpeningStateReq.ToState,
		changeOpeningStateReq.FromState,
		resultNoOpening,
		resultStateMismatch,
		resultUpdated,
	).Scan(&result)
	if err != nil {
		p.log.Err("failed to change opening state", "error", err)
		return db.ErrInternal
	}

	p.log.Dbg("state change result", "result", result)

	switch result {
	case resultNoOpening:
		return db.ErrNoOpening
	case resultStateMismatch:
		return db.ErrStateMismatch
	case resultUpdated:
		return nil
	default:
		p.log.Err("unexpected result from state change", "result", result)
		return db.ErrInternal
	}
}
