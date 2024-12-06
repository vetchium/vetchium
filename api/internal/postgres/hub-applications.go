package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) MyApplications(
	ctx context.Context,
	myApplicationsReq vetchi.MyApplicationsRequest,
) ([]vetchi.HubApplication, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("no hub user in context", "error", db.ErrNoHubUser)
		return []vetchi.HubApplication{}, db.ErrNoHubUser
	}

	p.log.Dbg("my applications request",
		"hubUserID", hubUser.ID,
		"state", myApplicationsReq.State,
		"paginationKey", myApplicationsReq.PaginationKey,
		"limit", myApplicationsReq.Limit)

	// First verify the hub user exists
	var userExists bool
	err := p.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM hub_users WHERE id = $1)
	`, hubUser.ID).Scan(&userExists)
	if err != nil {
		p.log.Err("failed to check hub user", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("hub user exists check", "exists", userExists)

	// Then check applications count
	var count int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications
		WHERE hub_user_id = $1
	`, hubUser.ID).Scan(&count)
	if err != nil {
		p.log.Err("failed to count applications", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("found applications count", "count", count)

	// First check the join between applications and openings
	var joinCount int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications a
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		WHERE a.hub_user_id = $1
	`, hubUser.ID).Scan(&joinCount)
	if err != nil {
		p.log.Err("failed to check join count", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("applications-openings join count", "count", joinCount)

	// Then check the full join path
	var fullJoinCount int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications a
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		JOIN employers e ON o.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		WHERE a.hub_user_id = $1
	`, hubUser.ID).Scan(&fullJoinCount)
	if err != nil {
		p.log.Err("failed to check full join count", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("full join count", "count", fullJoinCount)

	// First check just the hub_user_id condition
	var whereCount int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications a
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		JOIN employers e ON o.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		WHERE a.hub_user_id = $1
	`, hubUser.ID).Scan(&whereCount)
	if err != nil {
		p.log.Err("failed to check where count", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("where count (just hub_user_id)", "count", whereCount)

	// Convert empty state to nil for the database query
	var stateParam interface{}
	if myApplicationsReq.State == "" {
		stateParam = nil
	} else {
		stateParam = myApplicationsReq.State
	}

	// Convert pagination key to empty string if null
	var paginationParam string
	if myApplicationsReq.PaginationKey == nil {
		paginationParam = ""
	} else {
		paginationParam = *myApplicationsReq.PaginationKey
	}

	// Then check with the application state condition
	var stateCount int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications a
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		JOIN employers e ON o.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		WHERE a.hub_user_id = $1
		AND ($2::application_states IS NULL OR a.application_state = $2::application_states)
	`, hubUser.ID, stateParam).Scan(&stateCount)
	if err != nil {
		p.log.Err("failed to check state count", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("where count (with state)", "count", stateCount)

	// Finally check with all conditions
	var finalCount int
	err = p.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM applications a
		JOIN openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
		JOIN employers e ON o.employer_id = e.id
		JOIN employer_primary_domains epd ON e.id = epd.employer_id
		JOIN domains d ON epd.domain_id = d.id
		WHERE a.hub_user_id = $1
		AND ($2::application_states IS NULL OR a.application_state = $2::application_states)
		AND ($3 = '' OR a.id > $3)
	`, hubUser.ID, stateParam, paginationParam).Scan(&finalCount)
	if err != nil {
		p.log.Err("failed to check final count", "error", err)
		return []vetchi.HubApplication{}, err
	}
	p.log.Dbg("where count (all conditions)", "count", finalCount)

	query := `
SELECT
    a.id,
    a.application_state,
    a.opening_id,
    o.title,
    e.company_name,
    d.domain_name as employer_domain,
    a.created_at
FROM
    applications a
JOIN
    openings o ON a.employer_id = o.employer_id AND a.opening_id = o.id
JOIN
    employers e ON o.employer_id = e.id
JOIN
    employer_primary_domains epd ON e.id = epd.employer_id
JOIN
    domains d ON epd.domain_id = d.id
WHERE
    a.hub_user_id = $1
    AND ($2::application_states IS NULL OR a.application_state = $2::application_states)
    AND ($3 = '' OR a.id > $3)
ORDER BY
    a.created_at DESC,
    a.id ASC
LIMIT $4;
`

	p.log.Dbg("executing query",
		"query", query,
		"hubUserID", hubUser.ID,
		"state", myApplicationsReq.State,
		"paginationKey", myApplicationsReq.PaginationKey,
		"limit", myApplicationsReq.Limit)

	var hubApplications []vetchi.HubApplication
	rows, err := p.pool.Query(
		ctx,
		query,
		hubUser.ID,
		stateParam,
		paginationParam,
		myApplicationsReq.Limit,
	)
	if err != nil {
		p.log.Err("failed to get my applications", "error", err)
		return []vetchi.HubApplication{}, err
	}
	defer rows.Close()

	var rowCount int
	for rows.Next() {
		rowCount++
		var hubApplication vetchi.HubApplication
		if err := rows.Scan(
			&hubApplication.ApplicationID,
			&hubApplication.State,
			&hubApplication.OpeningID,
			&hubApplication.OpeningTitle,
			&hubApplication.EmployerName,
			&hubApplication.EmployerDomain,
			&hubApplication.CreatedAt,
		); err != nil {
			p.log.Err("failed to scan my applications", "error", err)
			return []vetchi.HubApplication{}, err
		}
		hubApplications = append(hubApplications, hubApplication)
	}

	p.log.Dbg("my applications result",
		"rowCount", rowCount,
		"hubApplications", hubApplications)
	return hubApplications, nil
}

func (p *PG) WithdrawApplication(
	ctx context.Context,
	applicationID string,
) error {
	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
		statusUpdated    = "updated"
	)

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("no hub user in context", "error", db.ErrNoHubUser)
		return db.ErrInternal
	}

	query := `
	WITH application_check AS (
		SELECT
			CASE
				WHEN NOT EXISTS (
					SELECT 1 FROM applications
					WHERE id = $1 AND hub_user_id = $2
				) THEN $4::text
				WHEN EXISTS (
					SELECT 1 FROM applications
					WHERE id = $1 AND hub_user_id = $2
					AND application_state != $3
				) THEN $5::text
				ELSE $6::text
			END as status
	), state_update AS (
		UPDATE applications a
		SET application_state = $7
		WHERE id = $1
		AND hub_user_id = $2
		AND application_state = $3
		AND EXISTS (SELECT 1 FROM application_check WHERE status = $6)
		RETURNING $8::text
	)
	SELECT COALESCE(
		(SELECT $8::text FROM state_update LIMIT 1),
		(SELECT status FROM application_check)
	);`

	var status string
	err := p.pool.QueryRow(
		ctx,
		query,
		applicationID,
		hubUser.ID,
		vetchi.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
		vetchi.WithdrawnAppState,
		statusUpdated,
	).Scan(&status)
	if err != nil {
		p.log.Err("failed to withdraw application", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		p.log.Dbg("application not found", "id", applicationID)
		return db.ErrNoApplication
	case statusWrongState:
		p.log.Dbg("application is in wrong state", "id", applicationID)
		return db.ErrApplicationStateInCompatible
	case statusUpdated:
		p.log.Dbg("withdrew application", "id", applicationID)
		return nil
	default:
		p.log.Err("unexpected status", "status", status)
		return db.ErrInternal
	}
}
