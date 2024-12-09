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

func (p *PG) AddLocation(
	ctx context.Context,
	addLocationRequest employer.AddLocationRequest,
) (uuid.UUID, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return uuid.UUID{}, db.ErrInternal
	}

	query := `
INSERT INTO locations (title, country_code, postal_address, postal_code, openstreetmap_url, city_aka, employer_id, location_state)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    id
`

	var locationID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		addLocationRequest.Title,
		addLocationRequest.CountryCode,
		addLocationRequest.PostalAddress,
		addLocationRequest.PostalCode,
		addLocationRequest.OpenStreetMapURL,
		addLocationRequest.CityAka,
		orgUser.EmployerID,
		employer.ActiveLocation,
	).Scan(&locationID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_location_title_employer_id" {
			return uuid.UUID{}, db.ErrDupLocationName
		}

		p.log.Err("failed to create location", "error", err)
		return uuid.UUID{}, err
	}

	return locationID, nil
}

func (p *PG) DefunctLocation(
	ctx context.Context,
	defunctLocationReq employer.DefunctLocationRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    locations
SET
    location_state = $1
WHERE
    title = $2
    AND employer_id = $3
RETURNING
    id
`

	_, err := p.pool.Exec(
		ctx,
		query,
		employer.DefunctLocation,
		defunctLocationReq.Title,
		orgUser.EmployerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoLocation
		}

		p.log.Err("failed to defunct location", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetLocByName(
	ctx context.Context,
	getLocationReq employer.GetLocationRequest,
) (employer.Location, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return employer.Location{}, db.ErrInternal
	}

	query := `
SELECT
    title,
    country_code,
    postal_address,
    postal_code,
    openstreetmap_url,
    city_aka,
    location_state
FROM
    locations
WHERE
    title = $1
    AND employer_id = $2
`

	var location employer.Location
	err := p.pool.QueryRow(
		ctx,
		query,
		getLocationReq.Title,
		orgUser.EmployerID,
	).Scan(
		&location.Title,
		&location.CountryCode,
		&location.PostalAddress,
		&location.PostalCode,
		&location.OpenStreetMapURL,
		&location.CityAka,
		&location.State,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return employer.Location{}, db.ErrNoLocation
		}

		p.log.Err("failed to get location by name", "error", err)
		return employer.Location{}, err
	}

	return location, nil
}

func (p *PG) GetLocations(
	ctx context.Context,
	getLocationsReq employer.GetLocationsRequest,
) ([]employer.Location, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return nil, db.ErrInternal
	}

	query := `
SELECT
    title,
    country_code,
    postal_address,
    postal_code,
    openstreetmap_url,
    city_aka,
    location_state
FROM
    locations
WHERE
    employer_id = $1
    AND location_state = ANY ($2::location_states[])
	AND title > $3
ORDER BY
    title ASC
LIMIT $4
`

	rows, err := p.pool.Query(
		ctx,
		query,
		orgUser.EmployerID,
		getLocationsReq.StatesAsStrings(),
		getLocationsReq.PaginationKey,
		getLocationsReq.Limit,
	)
	if err != nil {
		p.log.Err("failed to get locations", "error", err)
		return nil, err
	}

	locations, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[employer.Location],
	)
	if err != nil {
		p.log.Err("failed to get locations", "error", err)
		return nil, err
	}

	return locations, nil
}

func (p *PG) RenameLocation(
	ctx context.Context,
	renameLocationReq employer.RenameLocationRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    locations
SET
    title = $1
WHERE
    title = $2
    AND employer_id = $3
RETURNING
    id
`

	var locationID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		renameLocationReq.NewTitle,
		renameLocationReq.OldTitle,
		orgUser.EmployerID,
	).Scan(&locationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoLocation
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_location_name_employer_id" {
			return db.ErrDupLocationName
		}

		p.log.Err("failed to rename location", "error", err)
		return err
	}

	p.log.Dbg("location renamed", "location_id", locationID)

	return nil
}

func (p *PG) UpdateLocation(
	ctx context.Context,
	updateLocationReq employer.UpdateLocationRequest,
) error {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return db.ErrInternal
	}

	query := `
UPDATE
    locations
SET
    country_code = $1,
    postal_address = $2,
    postal_code = $3,
    openstreetmap_url = $4,
    city_aka = $5
WHERE
    title = $6
    AND employer_id = $7
RETURNING
    id
`

	var locationID uuid.UUID
	err := p.pool.QueryRow(
		ctx,
		query,
		updateLocationReq.CountryCode,
		updateLocationReq.PostalAddress,
		updateLocationReq.PostalCode,
		updateLocationReq.OpenStreetMapURL,
		updateLocationReq.CityAka,
		updateLocationReq.Title,
		orgUser.EmployerID,
	).Scan(&locationID)
	if err != nil {
		p.log.Err("failed to update location", "error", err)
		return err
	}

	p.log.Dbg("location updated", "location_id", locationID)

	return nil
}
