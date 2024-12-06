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
		return []vetchi.HubApplication{}, db.ErrNoHubUser
	}

	query := `
SELECT
    a.id,
    a.state,
    a.opening_id,
    o.title,
    o.employer_name,
    o.employer_domain,
    a.created_at
FROM
    hub_applications a
JOIN
    openings o ON a.opening_id = o.id
WHERE
    a.hub_user_id = $1
    AND (COALESCE($2, a.state) = a.state)
	AND a.created_at >= COALESCE(
		(SELECT created_at FROM hub_applications WHERE id = $3),
		'1970-01-01'
	)
	AND a.id > $3
ORDER BY
    a.created_at DESC,
    a.id ASC
LIMIT $4
`

	var hubApplications []vetchi.HubApplication
	rows, err := p.pool.Query(
		ctx,
		query,
		hubUser.ID,
		myApplicationsReq.State,
		myApplicationsReq.PaginationKey,
		myApplicationsReq.Limit,
	)
	if err != nil {
		p.log.Err("failed to get my applications", "error", err)
		return []vetchi.HubApplication{}, err
	}
	defer rows.Close()

	for rows.Next() {
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

	p.log.Dbg("my applications", "hubApplications", hubApplications)
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
	)

	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		p.log.Err("no hub user in context", "error", db.ErrNoHubUser)
		return db.ErrInternal
	}

	query := `
WITH application_check AS (
	SELECT CASE
		WHEN NOT EXISTS (
			SELECT 1 FROM applications
			WHERE id = $1 AND hub_user_id = $2
		) THEN $4
		WHEN EXISTS (
			SELECT 1 FROM applications
			WHERE id = $1 AND hub_user_id = $2
			AND application_state != $3
		) THEN $5
		ELSE $6
	END as status
)
UPDATE applications
SET application_state = $7
WHERE id = $1
AND hub_user_id = $2
AND application_state = $3
AND (SELECT status FROM application_check) = $6
RETURNING (SELECT status FROM application_check);
`

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
	case statusOK:
		p.log.Dbg("withdrew application", "id", applicationID)
		return nil
	default:
		p.log.Err("unexpected status", "status", status)
		return db.ErrInternal
	}
}
