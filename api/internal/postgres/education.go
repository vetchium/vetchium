package postgres

import (
	"context"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/hub"
)

func (pg *PG) AddEducation(
	ctx context.Context,
	addEducationReq hub.AddEducationRequest,
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
	`, addEducationReq.InstituteDomain).Scan(&instituteID)
	if err != nil {
		pg.log.Err(
			"failed to get or create dummy institute",
			"error",
			err,
			"domain",
			addEducationReq.InstituteDomain,
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
	if addEducationReq.StartDate != nil {
		startDate = *addEducationReq.StartDate
	}

	if addEducationReq.EndDate != nil {
		endDate = *addEducationReq.EndDate
	}

	err = tx.QueryRow(
		ctx,
		query,
		hubUserID,
		instituteID,
		addEducationReq.Degree,
		startDate,
		endDate,
		addEducationReq.Description,
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
	deleteEducationReq hub.DeleteEducationRequest,
) error {
	hubUserID, err := getHubUserID(ctx)
	if err != nil {
		pg.log.Err("failed to get hub user ID", "error", err)
		return db.ErrInternal
	}

	res, err := pg.pool.Exec(ctx, `
		DELETE FROM education
		WHERE id = $1 AND hub_user_id = $2
	`, deleteEducationReq.EducationID, hubUserID)
	if err != nil {
		pg.log.Err("failed to delete education", "error", err)
		return db.ErrInternal
	}

	if res.RowsAffected() == 0 {
		pg.log.Dbg("not found", "education_id", deleteEducationReq.EducationID)
		return db.ErrNoEducation
	}

	return nil
}

func (pg *PG) ListEducation(
	ctx context.Context,
	listEducationReq hub.ListEducationRequest,
) ([]hub.Education, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		return []hub.Education{}, db.ErrInternal
	}

	educationQuery := `
SELECT id, institute_domains.domain, degree, start_date, end_date, description
FROM education
JOIN institute_domains ON education.institute_id = institute_domains.id
WHERE hub_user_id = (
	SELECT id FROM hub_users WHERE handle = COALESCE($1, $2)
)
`

	rows, err := pg.pool.Query(
		ctx,
		educationQuery,
		listEducationReq.UserHandle,
		hubUser.Handle,
	)
	if err != nil {
		pg.log.Err("failed to query education", "error", err)
		return nil, err
	}

	educations := []hub.Education{}
	for rows.Next() {
		var education hub.Education
		err = rows.Scan(
			&education.ID,
			&education.InstituteDomain,
			&education.Degree,
			&education.StartDate,
			&education.EndDate,
			&education.Description,
		)
		if err != nil {
			pg.log.Err("failed to scan education", "error", err)
			return nil, err
		}

		if listEducationReq.UserHandle != nil &&
			*listEducationReq.UserHandle != hubUser.Handle {
			// Expose the ID only to the owner of the education
			education.ID = ""
		}

		educations = append(educations, education)
	}

	return educations, nil
}

func (pg *PG) FilterInstitutes(
	ctx context.Context,
	req hub.FilterInstitutesRequest,
) ([]hub.Institute, error) {
	return nil, nil
}
