package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/typespec/common"
)

// GetActiveApplicationScoringModels returns all active scoring models
func (p *PG) GetActiveApplicationScoringModels(
	ctx context.Context,
) ([]db.ApplicationScoringModel, error) {
	p.log.Dbg("getting active application scoring models")
	query := `
SELECT model_name, description, is_active, created_at
FROM application_scoring_models
WHERE is_active = true
`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		p.log.Err("Failed to query application scoring models", "error", err)
		return nil, err
	}
	defer rows.Close()

	var models []db.ApplicationScoringModel
	for rows.Next() {
		var model db.ApplicationScoringModel
		if err := rows.Scan(
			&model.ModelName,
			&model.Description,
			&model.IsActive,
			&model.CreatedAt,
		); err != nil {
			p.log.Err("Failed to scan application scoring model", "error", err)
			return nil, fmt.Errorf(
				"failed to scan application scoring model: %w",
				err,
			)
		}
		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("Error iterating application scoring models", "error", err)
		return nil, err
	}

	p.log.Dbg("Got active application scoring models", "count", len(models))
	return models, nil
}

// GetOpeningsWithUnscoredApplications returns openings that have applications in APPLIED state without scores
func (p *PG) GetOpeningsWithUnscoredApplications(
	ctx context.Context,
) ([]db.OpeningForScoring, error) {
	p.log.Dbg("getting openings with unscored applications")
	query := `
SELECT DISTINCT o.employer_id, o.id
FROM openings o
JOIN applications a ON o.employer_id = a.employer_id AND o.id = a.opening_id
WHERE a.application_state = $1
AND NOT EXISTS (
	SELECT 1 FROM application_scores s
	JOIN application_scoring_models m ON s.model_name = m.model_name
	WHERE s.application_id = a.id AND m.is_active = true
)
`

	rows, err := p.pool.Query(ctx, query, common.AppliedAppState)
	if err != nil {
		p.log.Err("query openings with unscored applications", "error", err)
		return nil, err
	}
	defer rows.Close()

	var openings []db.OpeningForScoring
	for rows.Next() {
		var opening db.OpeningForScoring
		if err := rows.Scan(&opening.EmployerID, &opening.ID); err != nil {
			p.log.Err("Failed to scan opening", "error", err)
			return nil, fmt.Errorf("failed to scan opening: %w", err)
		}
		openings = append(openings, opening)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("Error iterating openings", "error", err)
		return nil, fmt.Errorf("error iterating openings: %w", err)
	}

	p.log.Dbg("Got openings with unscored applications", "count", len(openings))
	return openings, nil
}

// GetOpeningJD returns the job description for an opening
func (p *PG) GetOpeningJD(
	ctx context.Context,
	employerID, openingID string,
) (string, error) {
	p.log.Dbg("Get JD", "employer_id", employerID, "opening_id", openingID)
	query := `
SELECT jd
FROM openings
WHERE employer_id = $1 AND id = $2
`

	var jd string
	err := p.pool.QueryRow(ctx, query, employerID, openingID).Scan(&jd)
	if err != nil {
		p.log.Err("JD", "emp", employerID, "opening", openingID, "err", err)
		return "", fmt.Errorf("failed to get opening JD: %w", err)
	}

	p.log.Dbg("Got JD", "employer_id", employerID, "opening_id", openingID)
	return jd, nil
}

// GetUnscoredApplicationsForOpening returns applications for an opening that have not been scored yet
func (p *PG) GetUnscoredApplicationsForOpening(
	ctx context.Context,
	employerID, openingID string,
	limit int,
) ([]db.ApplicationForScoring, error) {
	p.log.Dbg("unscored applications", "emp", employerID, "opening", openingID)

	query := `
SELECT a.id, a.resume_sha
FROM applications a
WHERE a.employer_id = $1
AND a.opening_id = $2
AND a.application_state = $3
AND NOT EXISTS (
	SELECT 1 FROM application_scores s
	JOIN application_scoring_models m ON s.model_name = m.model_name
	WHERE s.application_id = a.id AND m.is_active = true
)
LIMIT $4
`

	rows, err := p.pool.Query(
		ctx,
		query,
		employerID,
		openingID,
		common.AppliedAppState,
		limit,
	)
	if err != nil {
		p.log.Err("failed to query unscored applications", "err", err)
		return nil, fmt.Errorf("failed to query unscored applications: %w", err)
	}
	defer rows.Close()

	var applications []db.ApplicationForScoring
	for rows.Next() {
		var app db.ApplicationForScoring
		if err := rows.Scan(&app.ID, &app.ResumeSHA); err != nil {
			p.log.Err("Failed to scan application", "error", err)
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		p.log.Err("Error iterating applications", "error", err)
		return nil, fmt.Errorf("error iterating applications: %w", err)
	}

	p.log.Dbg("got unscored applications", "count", len(applications))
	return applications, nil
}

// SaveApplicationScores saves multiple scores for an application in a single transaction
func (p *PG) SaveApplicationScores(
	ctx context.Context,
	scores []db.ApplicationScore,
) error {
	p.log.Dbg("saving application scores", "count", len(scores))
	if len(scores) == 0 {
		p.log.Dbg("No scores to save, returning early")
		return nil
	}

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		p.log.Err("Failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Prepare the query
	query := `
INSERT INTO application_scores (application_id, model_name, score)
VALUES ($1, $2, $3)
ON CONFLICT (application_id, model_name) DO UPDATE
SET score = $3
`

	// Use a batch for efficiency
	batch := &pgx.Batch{}
	for _, score := range scores {
		batch.Queue(query, score.ApplicationID, score.ModelName, score.Score)
	}

	// Send the batch
	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	// Process each result
	for i := 0; i < batch.Len(); i++ {
		_, err = results.Exec()
		if err != nil {
			p.log.Err("Failed to execute batch query", "index", i, "error", err)
			return fmt.Errorf("failed to execute batch query: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(context.Background()); err != nil {
		p.log.Err("Failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	p.log.Dbg("Successfully saved batch of scores", "count", len(scores))
	return nil
}
