package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) SetApplicationColorTag(
	ctx context.Context,
	req employer.SetApplicationColorTagRequest,
) error {
	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
	)

	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
		WITH application_check AS (
			SELECT CASE
				WHEN NOT EXISTS (
					SELECT 1 FROM applications
					WHERE id = $2 AND employer_id = $3
				) THEN $5
				WHEN EXISTS (
					SELECT 1 FROM applications
					WHERE id = $2 AND employer_id = $3
					AND application_state != $4
				) THEN $6
				ELSE $7
			END as status
		),
		update_result AS (
			UPDATE applications
			SET color_tag = CASE 
				WHEN (SELECT status FROM application_check) = $7 THEN $1 
				ELSE color_tag 
			END
			WHERE id = $2 
			AND employer_id = $3
			AND application_state = $4
		)
		SELECT status FROM application_check;
	`

	var status string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.ColorTag,
		req.ApplicationID,
		orgUser.EmployerID,
		common.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to add application color tag", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		return nil
	default:
		p.log.Err("failed to add application color tag", "error", err)
		return db.ErrInternal
	}
}

func (p *PG) RemoveApplicationColorTag(
	ctx context.Context,
	req employer.RemoveApplicationColorTagRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	const (
		statusNotFound   = "not_found"
		statusWrongState = "wrong_state"
		statusOK         = "ok"
	)

	query := `
		WITH application_check AS (
			SELECT CASE
				WHEN NOT EXISTS (
					SELECT 1 FROM applications
					WHERE id = $1 AND employer_id = $2
				) THEN $4
				WHEN EXISTS (
					SELECT 1 FROM applications
					WHERE id = $1 AND employer_id = $2
					AND application_state != $3
				) THEN $5
				ELSE $6
			END as status
		),
		update_result AS (
			UPDATE applications
			SET color_tag = CASE 
				WHEN (SELECT status FROM application_check) = $6 THEN NULL 
				ELSE color_tag 
			END
			WHERE id = $1 
			AND employer_id = $2
			AND application_state = $3
		)
		SELECT status FROM application_check;
	`

	var status string
	err := p.pool.QueryRow(
		ctx,
		query,
		req.ApplicationID,
		orgUser.EmployerID,
		common.AppliedAppState,
		statusNotFound,
		statusWrongState,
		statusOK,
	).Scan(&status)

	if err != nil {
		p.log.Err("failed to remove application color tag", "error", err)
		return db.ErrInternal
	}

	switch status {
	case statusNotFound:
		return db.ErrNoApplication
	case statusWrongState:
		return db.ErrApplicationStateInCompatible
	case statusOK:
		return nil
	default:
		p.log.Err("unexpected status when removing color tag", "error", err)
		return db.ErrInternal
	}
}
