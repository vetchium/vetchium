package granger

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TimelineRefresher continuously refreshes user timelines
// It runs as a goroutine and stops when the quit channel is closed
func (g *Granger) TimelineRefresher(quit chan struct{}) {
	g.log.Dbg("TimelineRefresher started")

	// Wait time when no timelines need refreshing
	waitInterval := 1 * time.Minute

	// Standard timeout for database operations
	dbTimeout := 3 * time.Minute

	// Timer for waiting when no timelines are found
	timer := time.NewTimer(1 * time.Second) // Immediate first check
	defer timer.Stop()

	for {
		select {
		case <-quit:
			g.log.Dbg("TimelineRefresher stopped")
			return

		case <-timer.C:
			// Try to process a batch of timelines
			timelines, err := g.getTimelinesToRefresh(dbTimeout)
			if err != nil || len(timelines) == 0 {
				// Wait before checking again if error or no timelines
				timer.Reset(waitInterval)
				continue
			}

			// Process the timelines we found
			g.processTimelines(timelines, dbTimeout)

			// Check for more timelines immediately
			timer.Reset(1 * time.Second)
		}
	}
}

// getTimelinesToRefresh retrieves timelines that need to be refreshed
func (g *Granger) getTimelinesToRefresh(
	timeout time.Duration,
) ([]uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	g.log.Dbg("Checking for timelines to refresh")
	timelines, err := g.db.GetTimelinesToRefresh(ctx, 3)

	if err != nil {
		g.log.Err("Failed to get timelines to refresh", "error", err)
		return nil, err
	}

	if len(timelines) > 0 {
		g.log.Inf("Found timelines to refresh", "count", len(timelines))
	} else {
		g.log.Dbg("No timelines need refreshing")
	}

	return timelines, nil
}

// processTimelines refreshes a batch of timelines
func (g *Granger) processTimelines(
	timelines []uuid.UUID,
	timeout time.Duration,
) {
	for _, timeline := range timelines {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		g.log.Dbg("Refreshing timeline", "hub_user_id", timeline)
		err := g.db.RefreshTimeline(ctx, timeline)

		if err != nil {
			g.log.Err(
				"Failed to refresh timeline",
				"hub_user_id",
				timeline,
				"error",
				err,
			)
		} else {
			g.log.Dbg("Timeline refreshed successfully", "hub_user_id", timeline)
		}

		cancel()
	}
}
