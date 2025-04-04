package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/typespec/common"
)

// GetUnscoredApplication returns a random opening with unscored applications
func (p *PG) GetUnscoredApplication(
	ctx context.Context,
	limit int,
) (*db.UnscoredApplicationBatch, error) {
	query := `
WITH candidate_openings AS (
	SELECT DISTINCT o.employer_id, o.id, o.jd
	FROM openings o
	JOIN applications a ON o.employer_id = a.employer_id AND o.id = a.opening_id
	WHERE a.application_state = $1
	AND (o.state = $2 OR o.state = $3)
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
			return nil, nil
		}
		p.log.Err("Failed to query unscored application batch", "error", err)
		return nil, err
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

	// Execute individual SQL statements within the transaction
	for _, score := range scores {
		_, err := tx.Exec(
			ctx,
			query,
			score.ApplicationID,
			score.ModelName,
			score.Score,
		)
		if err != nil {
			p.log.Err("INSERT to application_scores failed", "error", err)
			return err
		}
	}

	// Commit the transaction
	if err = tx.Commit(context.Background()); err != nil {
		p.log.Err("Failed to commit transaction", "error", err)
		return err
	}

	p.log.Dbg("Successfully saved scores", "count", len(scores))
	return nil
}
