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

// GetUnscoredApplication returns a random opening with unscored applications
func (p *PG) GetUnscoredApplication(
	ctx context.Context,
	limit int,
) (*db.UnscoredApplicationBatch, error) {
	p.log.Dbg("getting unscored application batch")
	if limit <= 0 {
		limit = 10 // Default to max 10 applications
	}

	query := `
WITH candidate_openings AS (
	SELECT DISTINCT o.employer_id, o.id, o.jd
	FROM openings o
	JOIN applications a ON o.employer_id = a.employer_id AND o.id = a.opening_id
	WHERE a.application_state = $1
	AND (o.opening_state = $2 OR o.opening_state = $3)
	AND NOT EXISTS (
		SELECT 1 FROM application_scores s
		JOIN application_scoring_models m ON s.model_name = m.model_name
		WHERE s.application_id = a.id AND m.is_active = true
	)
	LIMIT 1
)
SELECT co.employer_id, co.id, co.jd, 
	array_agg(a.id) AS app_ids,
	array_agg(a.resume_sha) AS resume_shas
FROM candidate_openings co
JOIN applications a ON co.employer_id = a.employer_id AND co.id = a.opening_id
WHERE a.application_state = $1
AND NOT EXISTS (
	SELECT 1 FROM application_scores s
	JOIN application_scoring_models m ON s.model_name = m.model_name
	WHERE s.application_id = a.id AND m.is_active = true
)
GROUP BY co.employer_id, co.id, co.jd
LIMIT 1
`

	var employerID, openingID, jd string
	var appIDs, resumeSHAs []string

	err := p.pool.QueryRow(
		ctx,
		query,
		common.AppliedAppState,
		common.ActiveOpening,
		common.SuspendedOpening,
	).Scan(&employerID, &openingID, &jd, &appIDs, &resumeSHAs)

	if err != nil {
		if err == pgx.ErrNoRows {
			p.log.Dbg("No unscored applications found")
			return nil, nil
		}
		p.log.Err("Failed to query unscored application batch", "error", err)
		return nil, fmt.Errorf(
			"failed to query unscored application batch: %w",
			err,
		)
	}

	// Create result struct
	batch := &db.UnscoredApplicationBatch{
		EmployerID:   employerID,
		OpeningID:    openingID,
		JD:           jd,
		Applications: make([]db.ApplicationForScoring, 0, len(appIDs)),
	}

	// Populate applications (up to limit)
	maxApps := len(appIDs)
	if maxApps > limit {
		maxApps = limit
	}

	for i := 0; i < maxApps; i++ {
		batch.Applications = append(
			batch.Applications,
			db.ApplicationForScoring{
				ApplicationID: appIDs[i],
				ResumeSHA:     resumeSHAs[i],
			},
		)
	}

	p.log.Dbg("Got unscored application batch",
		"employer_id", employerID,
		"opening_id", openingID,
		"app_count", len(batch.Applications))

	return batch, nil
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
