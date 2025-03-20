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

	ticker := time.NewTicker(vetchi.ScoreApplicationsInterval)

	for {
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("Resume scoring job received quit signal")
			return
		case <-ticker.C:
			ticker.Stop()
			g.log.Dbg("Executing scoreApplications job")
			if err := g.processApplicationsForScoring(context.Background()); err != nil {
				g.log.Dbg("process applications for scoring", "error", err)
			}
			ticker = time.NewTicker(vetchi.ScoreApplicationsInterval)
		}
	}
}

func (g *Granger) processApplicationsForScoring(ctx context.Context) error {
	// Get openings with unscored applications in APPLIED state
	openings, err := g.db.GetOpeningsWithUnscoredApplications(ctx)
	if err != nil {
		return err
	}
	g.log.Dbg("Got openings with unscored applications", "count", len(openings))

	// Process each opening
	for _, opening := range openings {
		// Get job description for this opening
		jd, err := g.db.GetOpeningJD(ctx, opening.EmployerID, opening.ID)
		if err != nil {
			g.log.Dbg("failed to get job description", "err", err)
			continue
		}
		g.log.Dbg("Got job description", "jd_length", len(jd))

		// Get unscore applications for this opening (max 10 at a time)
		applications, err := g.db.GetUnscoredApplicationsForOpening(
			ctx,
			opening.EmployerID,
			opening.ID,
			vetchi.MaxApplicationsToScorePerBatch,
		)
		if err != nil {
			g.log.Dbg("failed to get unscored applications", "err", err)
			continue
		}
		g.log.Dbg("Got unscored applications", "count", len(applications))

		if len(applications) == 0 {
			continue
		}

		// Process this batch of applications
		err = g.scoreApplicationBatch(ctx, applications, jd)
		if err != nil {
			g.log.Dbg("failed to score application batch", "err", err)
			continue
		}
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
	resumePaths := make([]string, 0, len(applications))
	appIDMap := make(map[string]string) // Map resume paths to application IDs

	for _, app := range applications {
		// Format fileurl as expected by sortinghat: s3://bucket/key
		fileurl := fmt.Sprintf("s3://%s/%s", bucket, app.ResumeSHA)
		resumePaths = append(resumePaths, fileurl)
		appIDMap[fileurl] = app.ID
	}

	g.log.Dbg("Scoring batch of resumes", "count", len(resumePaths))

	// Create request payload
	request := vetchi.SortingHatRequest{
		JobDescription: jd,
		ResumePaths:    resumePaths,
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
