package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) GetCCByName(
	ctx context.Context,
	getCCByNameReq db.GetCCByNameReq,
) (vetchi.CostCenter, error) {
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
	var costCenter vetchi.CostCenter
	err := p.pool.QueryRow(ctx, query,
		getCCByNameReq.Name,
		getCCByNameReq.EmployerID,
	).Scan(&costCenter.Name,
		&costCenter.State,
		&costCenter.Notes,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return vetchi.CostCenter{}, db.ErrNoCostCenter
		}

		p.log.Error("failed to get cost center by name", "error", err)
		return vetchi.CostCenter{}, err
	}

	return costCenter, nil
}

func (p *PG) UpdateCostCenter(
	ctx context.Context,
	updateCCReq db.UpdateCCReq,
) error {
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
		updateCCReq.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoCostCenter
		}

		p.log.Error("failed to update cost center", "error", err)
		return err
	}

	p.log.Debug("cost center updated", "cost_center_id", costCenterID)

	return nil
}

func (p *PG) RenameCostCenter(
	ctx context.Context,
	renameCCReq db.RenameCCReq,
) error {
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
		renameCCReq.NewName,
		renameCCReq.OldName,
		renameCCReq.EmployerID,
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

		p.log.Error("failed to rename cost center", "error", err)
		return err
	}

	p.log.Debug("cost center renamed", "cost_center_id", costCenterID)
	return nil
}

func (p *PG) GetCostCenters(
	ctx context.Context,
	costCentersList db.CCentersList,
) ([]vetchi.CostCenter, error) {
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
		costCentersList.EmployerID,
		costCentersList.States,
		costCentersList.PaginationKey,
		costCentersList.Limit,
	)
	if err != nil {
		p.log.Error("failed to query cost centers", "error", err)
		return nil, err
	}

	costCenters, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[vetchi.CostCenter],
	)
	if err != nil {
		p.log.Error("failed to query cost centers", "error", err)
		return nil, err
	}

	return costCenters, nil
}

func (p *PG) DefunctCostCenter(
	ctx context.Context,
	defunctReq db.DefunctCCReq,
) error {
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
		vetchi.DefunctCC,
		defunctReq.Name,
		defunctReq.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoCostCenter
		}

		p.log.Error("failed to defunct cost center", "error", err)
		return err
	}

	return nil
}

func (p *PG) CreateCostCenter(
	ctx context.Context,
	costCenterReq db.CCenterReq,
) (uuid.UUID, error) {
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
		vetchi.ActiveCC,
		costCenterReq.Notes,
		costCenterReq.EmployerID,
	).Scan(&costCenterID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_cost_center_name_employer_id" {
			return uuid.UUID{}, db.ErrDupCostCenterName
		}

		p.log.Error("failed to create cost center", "error", err)
		return uuid.UUID{}, err
	}

	return costCenterID, nil
}
