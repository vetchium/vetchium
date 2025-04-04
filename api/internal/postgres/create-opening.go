package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

// TODO: I (PSankar) am not quite satisfied with the code in this file
// and this needs some cleanup with proper small functions, proper transactions,
// better error handling, proper retries, etc.

func (p *PG) CreateOpening(
	ctx context.Context,
	createOpeningReq employer.CreateOpeningRequest,
) (string, error) {
	orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
	if !ok {
		p.log.Err("failed to get orgUser from context")
		return "", db.ErrInternal
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(context.Background())

	var costCenterID uuid.UUID
	ccQuery := `
SELECT id FROM org_cost_centers WHERE cost_center_name = $1 AND employer_id = $2
`
	err = tx.QueryRow(ctx, ccQuery, createOpeningReq.CostCenterName, orgUser.EmployerID).
		Scan(&costCenterID)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("CC not found", "name", createOpeningReq.CostCenterName)
			return "", db.ErrNoCostCenter
		}
		p.log.Err("failed to get cost center", "error", err)
		return "", err
	}

	todayOpeningsCountQuery := `SELECT COUNT(*) FROM openings WHERE employer_id = $1 AND created_at::date = CURRENT_DATE`
	var todayOpeningsCount int
	err = tx.QueryRow(ctx, todayOpeningsCountQuery, orgUser.EmployerID).
		Scan(&todayOpeningsCount)
	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("first opening today", "employerID", orgUser.EmployerID)
			todayOpeningsCount = 0
		} else {
			p.log.Err("failed to get today's openings count", "error", err)
			return "", err
		}
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
	p.log.Dbg("generated opening ID", "id", openingID)

	var salaryMin, salaryMax *float64
	var currency *string
	if createOpeningReq.Salary != nil {
		salaryMin = &createOpeningReq.Salary.MinAmount
		salaryMax = &createOpeningReq.Salary.MaxAmount
		salaryCurrency := string(createOpeningReq.Salary.Currency)
		currency = &salaryCurrency
	}

	query := `
INSERT INTO openings (id, title, positions, jd, recruiter, hiring_manager, cost_center_id, employer_notes, remote_country_codes, remote_timezones, opening_type, yoe_min, yoe_max, min_education_level, salary_min, salary_max, salary_currency, state, employer_id)
    VALUES ($1, $2, $3, $4, (
            SELECT
                id
            FROM
                org_users
            WHERE
                email = $5
                AND employer_id = $19),
            (
                SELECT
                    id
                FROM
                    org_users
                WHERE
                    email = $6
                    AND employer_id = $19),
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
                $19)
    RETURNING
        id
`
	err = tx.QueryRow(
		ctx,
		query,
		openingID,
		createOpeningReq.Title,
		createOpeningReq.Positions,
		createOpeningReq.JD,
		createOpeningReq.Recruiter,
		createOpeningReq.HiringManager,
		costCenterID,
		createOpeningReq.EmployerNotes,
		createOpeningReq.RemoteCountryCodes,
		createOpeningReq.RemoteTimezones,
		createOpeningReq.OpeningType,
		createOpeningReq.YoeMin,
		createOpeningReq.YoeMax,
		createOpeningReq.MinEducationLevel,
		salaryMin,
		salaryMax,
		currency,
		common.DraftOpening,
		orgUser.EmployerID,
	).Scan(&openingID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23502 is not null error code
			if pgErr.Code == "23502" {
				if pgErr.ColumnName == "recruiter" {
					return "", db.ErrNoRecruiter
				} else if pgErr.ColumnName == "hiring_manager" {
					return "", db.ErrNoHiringManager
				} else {
					p.log.Err("create opening", "error", pgErr.Message)
					return "", err
				}
			}

			// Check for duplicate opening-id due to
			// race condition on parallel requests
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "openings_pkey" {
					// TODO: Handle this more gracefully
					p.log.Err("duplicate opening ID", "id", openingID)
					return "", db.ErrInternal
				}
			}

			p.log.Err("failed to create opening", "error", pgErr.Message)
			return "", err
		} else {
			p.log.Err("failed to create opening", "error", err)
		}
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
			if err == pgx.ErrNoRows {
				p.log.Dbg("not found", "team", createOpeningReq.HiringTeam)
				return "", db.ErrInvalidHiringTeam
			}
			p.log.Err("failed to verify hiring team", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Dbg(
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
			p.log.Err("failed to insert hiring team", "error", err)
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
			if err == pgx.ErrNoRows {
				p.log.Dbg("", "locations", createOpeningReq.LocationTitles)
				return "", db.ErrNoLocation
			}
			p.log.Err("failed to verify locations", "error", err)
			return "", err
		}
		if invalidCount > 0 {
			p.log.Dbg("invalid locations found", "count", invalidCount)
			return "", db.ErrNoLocation
		}

		// If all locations are valid, proceed with insertion
		locationQuery := `
INSERT INTO opening_locations (employer_id, opening_id, location_id)
SELECT $1, $2, l.id
FROM locations l
WHERE l.title = ANY($3)
AND l.employer_id = $1
`
		_, err = tx.Exec(
			ctx,
			locationQuery,
			orgUser.EmployerID,
			openingID,
			createOpeningReq.LocationTitles,
		)
		if err != nil {
			p.log.Err("failed to insert locations", "error", err)
			return "", err
		}
	}

	// Handle tags
	if len(createOpeningReq.Tags) > 0 {
		insertTagMappingsQuery := `
INSERT INTO opening_tag_mappings (employer_id, opening_id, tag_id)
SELECT $1, $2, UNNEST($3::uuid[])
`
		_, err = tx.Exec(
			ctx,
			insertTagMappingsQuery,
			orgUser.EmployerID,
			openingID,
			createOpeningReq.Tags,
		)
		if err != nil {
			p.log.Err("failed to insert tag mappings", "error", err)
			return "", err
		}
	}

	// Handle new tags
	if len(createOpeningReq.NewTags) > 0 {
		// First get IDs of any pre-existing tags from the names
		getExistingTagsQuery := `
SELECT id, name FROM opening_tags WHERE name = ANY($1)
`
		rows, err := tx.Query(
			ctx,
			getExistingTagsQuery,
			createOpeningReq.NewTags,
		)
		if err != nil {
			p.log.Err("failed to get existing tag ids", "error", err)
			return "", err
		}
		defer rows.Close()

		var tagIDs []uuid.UUID
		existingTagNames := make(map[string]struct{})
		for rows.Next() {
			var tagID uuid.UUID
			var tagName string
			err = rows.Scan(&tagID, &tagName)
			if err != nil {
				p.log.Err("failed to scan existing tag id", "error", err)
				return "", err
			}
			tagIDs = append(tagIDs, tagID)
			existingTagNames[tagName] = struct{}{}
		}

		// Find which tags need to be created
		var newTagNames []string
		for _, tagName := range createOpeningReq.NewTags {
			if _, exists := existingTagNames[tagName]; !exists {
				newTagNames = append(newTagNames, tagName)
			}
		}

		// Create only the tags that don't exist
		if len(newTagNames) > 0 {
			insertNewTagsQuery := `
INSERT INTO opening_tags (name)
SELECT UNNEST($1::text[])
RETURNING id
`
			rows, err := tx.Query(ctx, insertNewTagsQuery, newTagNames)
			if err != nil {
				p.log.Err("failed to insert new tags", "error", err)
				return "", err
			}
			defer rows.Close()

			for rows.Next() {
				var tagID uuid.UUID
				err = rows.Scan(&tagID)
				if err != nil {
					p.log.Err("failed to scan tag id", "error", err)
					return "", err
				}
				tagIDs = append(tagIDs, tagID)
			}
		}

		// Insert mappings for all tags
		insertTagMappingsQuery := `
INSERT INTO opening_tag_mappings (employer_id, opening_id, tag_id)
SELECT $1, $2, UNNEST($3::uuid[])
`
		_, err = tx.Exec(
			ctx,
			insertTagMappingsQuery,
			orgUser.EmployerID,
			openingID,
			tagIDs,
		)
		if err != nil {
			p.log.Err("failed to insert new tag mappings", "error", err)
			return "", err
		}
	}

	// Verify that at least one tag exists for the opening
	verifyTagsQuery := `
SELECT COUNT(*) FROM opening_tag_mappings
WHERE employer_id = $1 AND opening_id = $2
`
	var tagCount int
	err = tx.QueryRow(
		ctx,
		verifyTagsQuery,
		orgUser.EmployerID,
		openingID,
	).Scan(&tagCount)
	if err != nil {
		p.log.Err("failed to verify tag count", "error", err)
		return "", err
	}

	if tagCount == 0 {
		p.log.Dbg("no tags specified for opening")
		return "", errors.New("at least one tag is required for an opening")
	}

	if tagCount > 3 {
		p.log.Dbg("too many tags specified for opening")
		return "", errors.New("maximum of three tags allowed per opening")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		p.log.Err("failed to commit transaction", "error", err)
		return "", err
	}

	return openingID, nil
}
