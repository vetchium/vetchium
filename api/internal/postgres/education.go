package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
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
) ([]common.Education, error) {
	hubUser, ok := ctx.Value(middleware.HubUserCtxKey).(db.HubUserTO)
	if !ok {
		return []common.Education{}, db.ErrInternal
	}

	educationQuery := `
SELECT education.id, institute_domains.domain, education.degree, education.start_date, education.end_date, education.description
FROM education
JOIN institutes ON education.institute_id = institutes.id
JOIN institute_domains ON institutes.id = institute_domains.institute_id
WHERE education.hub_user_id = (
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
	defer rows.Close()

	educations := []common.Education{}
	for rows.Next() {
		var education common.Education
		var startDate, endDate interface{}
		var id uuid.UUID

		err = rows.Scan(
			&id,
			&education.InstituteDomain,
			&education.Degree,
			&startDate,
			&endDate,
			&education.Description,
		)
		if err != nil {
			pg.log.Err("failed to scan education", "error", err)
			return nil, err
		}

		// Convert UUID to string
		education.ID = id.String()

		// Convert date to string if not null
		if startDate != nil {
			startDateStr := startDate.(time.Time).Format("2006-01-02")
			education.StartDate = &startDateStr
		}

		if endDate != nil {
			endDateStr := endDate.(time.Time).Format("2006-01-02")
			education.EndDate = &endDateStr
		}

		if listEducationReq.UserHandle != nil &&
			string(*listEducationReq.UserHandle) != hubUser.Handle {
			// Expose the ID only to the owner of the education
			education.ID = ""
		}

		educations = append(educations, education)
	}

	if err = rows.Err(); err != nil {
		pg.log.Err("error iterating rows", "error", err)
		return nil, err
	}

	return educations, nil
}

func (pg *PG) FilterInstitutes(
	ctx context.Context,
	filterInstitutesReq hub.FilterInstitutesRequest,
) ([]common.Institute, error) {
	rows, err := pg.pool.Query(ctx, `
		SELECT id.domain, COALESCE(i.institute_name, id.domain) as name 
		FROM institute_domains id
		JOIN institutes i ON id.institute_id = i.id
		WHERE id.domain ILIKE $1 OR i.institute_name ILIKE $1
		LIMIT 10
	`, "%"+filterInstitutesReq.Prefix+"%")
	if err != nil {
		pg.log.Err("failed to filter institutes", "error", err)
		return nil, err
	}

	institutes := []common.Institute{}
	for rows.Next() {
		var institute common.Institute
		err = rows.Scan(&institute.Domain, &institute.Name)
		if err != nil {
			pg.log.Err("failed to scan institute", "error", err)
			return nil, err
		}

		institutes = append(institutes, institute)
	}

	return institutes, nil
}

func (pg *PG) ListHubUserEducation(
	ctx context.Context,
	listHubUserEducationReq employer.ListHubUserEducationRequest,
) ([]common.Education, error) {
	educationQuery := `
SELECT institute_domains.domain, education.degree, education.start_date, education.end_date, education.description
FROM education
JOIN institutes ON education.institute_id = institutes.id
JOIN institute_domains ON institutes.id = institute_domains.institute_id
WHERE education.hub_user_id = (
	SELECT id FROM hub_users WHERE handle = $1
)
`

	rows, err := pg.pool.Query(
		ctx,
		educationQuery,
		listHubUserEducationReq.Handle,
	)
	if err != nil {
		pg.log.Err("failed to query education", "error", err)
		return nil, err
	}
	defer rows.Close()

	educations := []common.Education{}
	for rows.Next() {
		var education common.Education
		var startDate, endDate interface{}

		err = rows.Scan(
			&education.InstituteDomain,
			&education.Degree,
			&startDate,
			&endDate,
			&education.Description,
		)
		if err != nil {
			pg.log.Err("failed to scan education", "error", err)
			return nil, err
		}

		// Convert date to string if not null
		if startDate != nil {
			startDateStr := startDate.(time.Time).Format("2006-01-02")
			education.StartDate = &startDateStr
		}

		if endDate != nil {
			endDateStr := endDate.(time.Time).Format("2006-01-02")
			education.EndDate = &endDateStr
		}

		educations = append(educations, education)
	}

	if err = rows.Err(); err != nil {
		pg.log.Err("error iterating rows", "error", err)
		return nil, err
	}

	return educations, nil
}
