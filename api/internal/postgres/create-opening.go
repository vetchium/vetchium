package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (p *PG) CreateOpening(
	ctx context.Context,
	createOpeningReq vetchi.CreateOpeningRequest,
) (string, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Error("failed to get orgUser from context")
		return "", db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var costCenterID uuid.UUID
	ccQuery := `SELECT id FROM cost_centers WHERE name = $1`
	err = tx.QueryRow(ctx, ccQuery, createOpeningReq.CostCenterName).
		Scan(&costCenterID)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Debug("CC not found", "name", createOpeningReq.CostCenterName)
			return "", db.ErrNoCostCenter
		}
		return "", err
	}

	todayOpeningsCountQuery := `SELECT COUNT(*) FROM openings WHERE employer_id = $1 AND created_at::date = CURRENT_DATE`
	var todayOpeningsCount int
	err = tx.QueryRow(ctx, todayOpeningsCountQuery, orgUser.EmployerID).
		Scan(&todayOpeningsCount)
	if err != nil {
		p.log.Error("failed to get today's openings count", "error", err)
		return "", err
	}

	// TODO: Check for max openings allowed per Employer and/or per Day.
	// Potentially this is where charging/pricing status codes will be returned
	var openingID string
	t := time.Now()
	openingID = fmt.Sprintf(
		"%d-%s-%d-%d",
		t.Year(),
		t.Format("Jan"),
		t.Day(),
		todayOpeningsCount+1,
	)
	p.log.Debug("generated opening ID", "id", openingID)

	query := `
INSERT INTO openings (id, title, positions, jd, recruiter, hiring_manager, cost_center_id, employer_notes, remote_country_codes, remote_timezones, opening_type, yoe_min, yoe_max, min_education_level, salary_min, salary_max, salary_currency, current_state, approval_waiting_state, employer_id)
    VALUES ($1, $2, $3, $4, (
            SELECT
                id
            FROM
                org_users
            WHERE
                email = $5
                AND employer_id = $21),
            (
                SELECT
                    id
                FROM
                    org_users
                WHERE
                    email = $6
                    AND employer_id = $21),
                $7,
                $8,
                $9,
                $10,
                $11,
                $12,
                $13,
                $14,
                $15,
                $16,
                $17,
                $18,
                $19,
                $20,
                $21)
    RETURNING
        id
`
	err = tx.QueryRow(ctx, query, openingID, createOpeningReq.Title, createOpeningReq.Positions, createOpeningReq.JD, createOpeningReq.HiringManager, costCenterID, createOpeningReq.EmployerNotes, createOpeningReq.RemoteCountryCodes, createOpeningReq.RemoteTimezones, createOpeningReq.OpeningType, createOpeningReq.YoeMin, createOpeningReq.YoeMax, createOpeningReq.MinEducationLevel, createOpeningReq.Salary.MinAmount, createOpeningReq.Salary.MaxAmount, createOpeningReq.Salary.Currency, vetchi.DraftOpening, nil, orgUser.EmployerID).
		Scan(&openingID)
	if err != nil {
		p.log.Error("failed to create opening", "error", err)
		return "", err
	}

	if len(createOpeningReq.HiringTeam) > 0 {
		// First verify all hiring team members exist and are in valid states
		verifyHiringTeamQuery := `
SELECT COUNT(*)
FROM (
    SELECT UNNEST($1::text[]) AS email
    EXCEPT
    SELECT email FROM org_users
    WHERE employer_id = $2
    AND org_user_state IN ('ACTIVE_ORG_USER', 'REPLICATED_ORG_USER')
) AS invalid_hiring_team`

		var invalidCount int
		err = tx.QueryRow(
			ctx,
			verifyHiringTeamQuery,
			createOpeningReq.HiringTeam,
			orgUser.EmployerID,
		).Scan(&invalidCount)
		if err != nil {
			p.log.Error("failed to verify hiring team", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Debug(
				"invalid hiring team members found",
				"count",
				invalidCount,
			)
			return "", db.ErrInvalidHiringTeam
		}

		// If all hiring team members are valid, proceed with insertion
		insertHiringTeamQuery := `
INSERT INTO opening_hiring_team (employer_id, opening_id, hiring_team_mate_id)
SELECT $1, $2, id
FROM org_users
WHERE email = ANY($3)
AND employer_id = $1
AND org_user_state IN ('ACTIVE_ORG_USER', 'REPLICATED_ORG_USER')
`
		_, err = tx.Exec(
			ctx,
			insertHiringTeamQuery,
			orgUser.EmployerID,
			openingID,
			createOpeningReq.HiringTeam,
		)
		if err != nil {
			p.log.Error("failed to insert hiring team", "error", err)
			return "", err
		}
	}

	if len(createOpeningReq.LocationTitles) > 0 {
		// First verify all locations exist for this employer
		verifyLocationsQuery := `
SELECT COUNT(*)
FROM (
    SELECT UNNEST($1::text[]) AS title
    EXCEPT
    SELECT title FROM locations
    WHERE employer_id = $2
) AS invalid_locations`

		var invalidCount int
		err = tx.QueryRow(
			ctx,
			verifyLocationsQuery,
			createOpeningReq.LocationTitles,
			orgUser.EmployerID,
		).Scan(&invalidCount)
		if err != nil {
			p.log.Error("failed to verify locations", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Debug("invalid locations found", "count", invalidCount)
			return "", db.ErrInvalidLocation
		}

		// If all locations are valid, proceed with insertion
		locationQuery := `
INSERT INTO opening_locations (opening_id, location_id)
SELECT $1, l.id
FROM locations l
WHERE l.title = ANY($2)
AND l.employer_id = $3
`
		_, err = tx.Exec(
			ctx,
			locationQuery,
			openingID,
			createOpeningReq.LocationTitles,
			orgUser.EmployerID,
		)
		if err != nil {
			p.log.Error("failed to insert locations", "error", err)
			return "", err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.log.Error("failed to commit transaction", "error", err)
		return "", err
	}

	return openingID, nil
}
