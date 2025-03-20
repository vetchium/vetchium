package granger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (g *Granger) scoreApplications(quit <-chan struct{}) {
	g.log.Dbg("Starting scoreApplications job")
	defer g.log.Dbg("scoreApplications job finished")
	defer g.wg.Done()

	for {
		ticker := time.NewTicker(vetchi.ScoreApplicationsInterval)
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("Resume scoring job received quit signal")
			return
		case <-ticker.C:
			ticker.Stop()
			if err := g.processApplicationsForScoring(context.Background()); err != nil {
				g.log.Dbg("process applications for scoring", "error", err)
			}
		}
	}
}

func (g *Granger) processApplicationsForScoring(ctx context.Context) error {
	// Get an unscored application batch
	batch, err := g.db.GetUnscoredApplication(
		ctx,
		vetchi.MaxApplicationsToScorePerBatch,
	)
	if err != nil {
		return err
	}

	// If no unscored applications found, return early
	if batch == nil {
		// No unscored applications found
		return nil
	}

	g.log.Dbg("Got unscored application batch",
		"employer_id", batch.EmployerID,
		"opening_id", batch.OpeningID,
		"app_count", len(batch.Applications))

	// Score the batch of applications
	err = g.scoreApplicationBatch(ctx, batch.Applications, batch.JD)
	if err != nil {
		g.log.Dbg("failed to score application batch", "err", err)
		return err
	}

	return nil
}

func (g *Granger) scoreApplicationBatch(
	ctx context.Context,
	applications []db.ApplicationForScoring,
	jd string,
) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Get S3 bucket name from environment variable
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		g.log.Err("S3_BUCKET environment variable not set")
		return fmt.Errorf("S3_BUCKET environment variable not set")
	}

	// Prepare batch request
	appSortRequests := make(
		[]vetchi.ApplicationSortRequest,
		0,
		len(applications),
	)
	appIDMap := make(map[string]string) // Map resume paths to application IDs

	for _, app := range applications {
		// Format fileurl as expected by sortinghat: s3://bucket/key
		fileurl := fmt.Sprintf("s3://%s/%s", bucket, app.ResumeSHA)
		appSortRequests = append(appSortRequests, vetchi.ApplicationSortRequest{
			ApplicationID: app.ApplicationID,
			ResumePath:    fileurl,
		})
		appIDMap[fileurl] = app.ApplicationID
	}

	g.log.Dbg("Scoring batch of resumes", "count", len(appSortRequests))

	// Create request payload
	request := vetchi.SortingHatRequest{
		JobDescription:          jd,
		ApplicationSortRequests: appSortRequests,
	}

	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		g.log.Err("failed to marshal request", "err", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build request
	apiURL := "http://sortinghat:8080/score-batch"
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		apiURL,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		g.log.Err("failed to create request", "err", err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		g.log.Err("failed to call sortinghat API", "err", err)
		return fmt.Errorf("failed to call sortinghat API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.log.Err(
			"sortinghat API returned non-OK status",
			"status",
			resp.Status,
		)
		return fmt.Errorf("sortinghat API returned status %s", resp.Status)
	}

	// Parse response
	var response vetchi.SortingHatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		g.log.Err("failed to decode sortinghat response", "err", err)
		return fmt.Errorf("failed to decode sortinghat response: %w", err)
	}

	g.log.Dbg("Sortinghat response", "response", response)

	// Collect all scores to save in a single transaction
	var allScores []db.ApplicationScore

	// Process scores from the response
	for _, score := range response.Scores {
		appID := score.ApplicationID

		// Map model scores to database scores
		for _, modelScore := range score.ModelScores {
			allScores = append(allScores, db.ApplicationScore{
				ApplicationID: appID,
				ModelName:     modelScore.ModelName,
				Score:         modelScore.Score,
			})
		}
	}

	// Save all scores in a single transaction
	if len(allScores) > 0 {
		err := g.db.SaveApplicationScores(ctx, allScores)
		if err != nil {
			g.log.Err("failed to save application scores", "err", err)
			return fmt.Errorf("failed to save application scores: %w", err)
		}
		g.log.Dbg("Saved all application scores", "count", len(allScores))
	} else {
		g.log.Dbg("No scores to save")
	}

	return nil
}
