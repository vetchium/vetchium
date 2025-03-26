package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/hub"
)

func (pg *PG) AddEducation(
	ctx context.Context,
	req hub.AddEducationRequest,
) (string, error) {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return "", err
	}

	tx, err := pg.pool.Begin(ctx)
	if err != nil {
		pg.log.Err("failed to begin transaction", "error", err)
		return "", err
	}
	defer tx.Rollback(context.Background())

	// Get or create institute for the domain
	var instituteID string
	err = tx.QueryRow(ctx, `
		SELECT get_or_create_dummy_institute($1)
	`, req.InstituteDomain).Scan(&instituteID)
	if err != nil {
		pg.log.Err(
			"failed to get or create dummy institute",
			"error",
			err,
			"domain",
			req.InstituteDomain,
		)
		return "", db.ErrInternal
	}

	// Insert education
	var id string
	query := `
INSERT INTO education (
	hub_user_id,
	institute_id,
	degree,
	start_date,
	end_date,
	description
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`
	// Handle possible null date values
	var startDate, endDate interface{}

	// Use nil for empty/invalid dates to properly handle NULL in database
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	err = tx.QueryRow(
		ctx,
		query,
		hubUserID,
		instituteID,
		req.Degree,
		startDate,
		endDate,
		req.Description,
	).Scan(&id)
	if err != nil {
		pg.log.Err("failed to insert education", "error", err)
		return "", db.ErrInternal
	}

	err = tx.Commit(context.Background())
	if err != nil {
		pg.log.Err("failed to commit transaction", "error", err)
		return "", db.ErrInternal
	}

	return id, nil
}

func (pg *PG) DeleteEducation(
	ctx context.Context,
	req hub.DeleteEducationRequest,
) error {
	return nil
}

func (pg *PG) ListEducation(
	ctx context.Context,
	req hub.ListEducationRequest,
) ([]hub.Education, error) {
	return nil, nil
}

func (pg *PG) FilterInstitutes(
	ctx context.Context,
	req hub.FilterInstitutesRequest,
) ([]hub.Institute, error) {
	return nil, nil
}
