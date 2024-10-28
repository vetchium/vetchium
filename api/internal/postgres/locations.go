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

func (p *PG) AddLocation(
	ctx context.Context,
	req db.AddLocationReq,
) (uuid.UUID, error) {
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
		req.Title,
		req.CountryCode,
		req.PostalAddress,
		req.PostalCode,
		req.OpenStreetMapURL,
		req.CityAka,
		req.EmployerID,
		vetchi.ActiveLocation,
	).Scan(&locationID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" &&
			pgErr.ConstraintName == "uniq_location_title_employer_id" {
			return uuid.UUID{}, db.ErrDupLocationName
		}

		p.log.Error("failed to create location", "error", err)
		return uuid.UUID{}, err
	}

	return locationID, nil
}

func (p *PG) DefunctLocation(
	ctx context.Context,
	req db.DefunctLocationReq,
) error {
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
		vetchi.DefunctLocation,
		req.Title,
		req.EmployerID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.ErrNoLocation
		}

		p.log.Error("failed to defunct location", "error", err)
		return err
	}

	return nil
}

func (p *PG) GetLocByName(
	ctx context.Context,
	req db.GetLocByNameReq,
) (vetchi.Location, error) {
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

	var location vetchi.Location
	err := p.pool.QueryRow(
		ctx,
		query,
		req.Title,
		req.EmployerID,
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
			return vetchi.Location{}, db.ErrNoLocation
		}

		p.log.Error("failed to get location by name", "error", err)
		return vetchi.Location{}, err
	}

	return location, nil
}

func (p *PG) GetLocations(
	ctx context.Context,
	req db.GetLocationsReq,
) ([]vetchi.Location, error) {
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
		req.EmployerID,
		req.States,
		req.PaginationKey,
		req.Limit,
	)
	if err != nil {
		p.log.Error("failed to get locations", "error", err)
		return nil, err
	}

	locations, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[vetchi.Location],
	)
	if err != nil {
		p.log.Error("failed to get locations", "error", err)
		return nil, err
	}

	return locations, nil
}

func (p *PG) RenameLocation(
	ctx context.Context,
	req db.RenameLocationReq,
) error {
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
		req.NewTitle,
		req.OldTitle,
		req.EmployerID,
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

		p.log.Error("failed to rename location", "error", err)
		return err
	}

	p.log.Debug("location renamed", "location_id", locationID)

	return nil
}

func (p *PG) UpdateLocation(
	ctx context.Context,
	req db.UpdateLocationReq,
) error {
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
		req.CountryCode,
		req.PostalAddress,
		req.PostalCode,
		req.OpenStreetMapURL,
		req.CityAka,
		req.Title,
		req.EmployerID,
	).Scan(&locationID)
	if err != nil {
		p.log.Error("failed to update location", "error", err)
		return err
	}

	p.log.Debug("location updated", "location_id", locationID)

	return nil
}
