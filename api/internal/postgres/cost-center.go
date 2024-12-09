package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/employer"
)

func (p *PG) GetCCByName(
	ctx context.Context,
	getCCByNameReq employer.GetCostCenterRequest,
) (employer.CostCenter, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.CostCenter{}, db.ErrInternal
	}

	query := `
SELECT
    cost_center_name,
    cost_center_state,
    notes
FROM
    org_cost_centers
WHERE
    cost_center_name = $1
    AND employer_id = $2
`

	// TODO: Perhaps in the future we will want to use sqlx.ScanStruct
	// but for now this is fine.
	var costCenter employer.CostCenter
	err := p.pool.QueryRow(ctx, query,
		getCCByNameReq.Name,
		orgUser.EmployerID,
	).Scan(&costCenter.Name,
		&costCenter.State,
		&costCenter.Notes,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return employer.CostCenter{}, db.ErrNoCostCenter
		}

		p.log.Err("failed to get cost center by name", "error", err)
		return employer.CostCenter{}, err
	}

	return costCenter, nil
}

func (p *PG) UpdateCostCenter(
	ctx context.Context,
	updateCCReq employer.UpdateCostCenterRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    org_cost_centers
SET
    notes = $1
WHERE
    cost_center_name = $2
    AND employer_id = $3
RETURNING id
`
	var costCenterID uuid.UUID
	err := p.pool.QueryRow(ctx, query,
		updateCCReq.Notes,
		updateCCReq.Name,
		orgUser.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoCostCenter
		}

		p.log.Err("failed to update cost center", "error", err)
		return err
	}

	p.log.Dbg("cost center updated", "cost_center_id", costCenterID)

	return nil
}

func (p *PG) RenameCostCenter(
	ctx context.Context,
	renameCostCenterReq employer.RenameCostCenterRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    org_cost_centers
SET
    cost_center_name = $1
WHERE
    cost_center_name = $2
    AND employer_id = $3
RETURNING id
`

	var costCenterID uuid.UUID
	err := p.pool.QueryRow(ctx, query,
		renameCostCenterReq.NewName,
		renameCostCenterReq.OldName,
		orgUser.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoCostCenter
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_cost_center_name_employer_id" {
			return db.ErrDupCostCenterName
		}

		p.log.Err("failed to rename cost center", "error", err)
		return err
	}

	p.log.Dbg("cost center renamed", "cost_center_id", costCenterID)
	return nil
}

func (p *PG) GetCostCenters(
	ctx context.Context,
	costCentersList employer.GetCostCentersRequest,
) ([]employer.CostCenter, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	query := `
SELECT
    oc.cost_center_name,
	oc.cost_center_state,
    oc.notes
FROM
    org_cost_centers oc
WHERE
    oc.employer_id = $1::uuid
    AND oc.cost_center_state = ANY ($2::cost_center_states[])
	AND oc.cost_center_name > $3
ORDER BY
    oc.cost_center_name ASC
LIMIT $4
`

	rows, err := p.pool.Query(ctx, query,
		orgUser.EmployerID,
		costCentersList.StatesAsStrings(),
		costCentersList.PaginationKey,
		costCentersList.Limit,
	)
	if err != nil {
		p.log.Err("failed to query cost centers", "error", err)
		return nil, err
	}

	costCenters, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[employer.CostCenter],
	)
	if err != nil {
		p.log.Err("failed to query cost centers", "error", err)
		return nil, err
	}

	return costCenters, nil
}

func (p *PG) DefunctCostCenter(
	ctx context.Context,
	defunctCostCenterReq employer.DefunctCostCenterRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    org_cost_centers
SET
    cost_center_state = $1
WHERE
    cost_center_name = $2
    AND employer_id = $3
RETURNING
    id
`

	var costCenterID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		employer.DefunctCC,
		defunctCostCenterReq.Name,
		orgUser.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoCostCenter
		}

		p.log.Err("failed to defunct cost center", "error", err)
		return err
	}

	return nil
}

func (p *PG) CreateCostCenter(
	ctx context.Context,
	costCenterReq employer.AddCostCenterRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
INSERT INTO org_cost_centers (cost_center_name, cost_center_state, notes, employer_id)
    VALUES ($1, $2, $3, $4)
RETURNING
    id
`
	var costCenterID uuid.UUID
	err := p.pool.QueryRow(
		ctx, query,
		costCenterReq.Name,
		employer.ActiveCC,
		costCenterReq.Notes,
		orgUser.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_cost_center_name_employer_id" {
			return uuid.UUID{}, db.ErrDupCostCenterName
		}

		p.log.Err("failed to create cost center", "error", err)
		return uuid.UUID{}, err
	}

	return costCenterID, nil
}
