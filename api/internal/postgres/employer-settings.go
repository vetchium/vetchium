package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
)

func (pg *PG) ChangeCoolOffPeriod(
	ctx context.Context,
	coolOffPeriod int32,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		pg.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	// TODO: Audit logs
	_, err := pg.pool.Exec(ctx, `
		UPDATE employers
		SET cool_off_period_days = $1
		WHERE id = $2
	`, coolOffPeriod, orgUser.EmployerID)
	if err != nil {
		pg.log.Err("failed to change cool off period", "error", err)
		return err
	}

	pg.log.Dbg("cool off period changed", "cool off period", coolOffPeriod)

	return nil
}

func (pg *PG) GetCoolOffPeriod(ctx context.Context) (int32, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		pg.log.Err("failed to get orgUser from context")
		return -1, db.ErrInternal
	}

	var period int32
	err := pg.pool.QueryRow(ctx, `
		SELECT cool_off_period_days FROM employers
		WHERE id = $1
	`, orgUser.EmployerID).Scan(&period)
	return period, err
}
