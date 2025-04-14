package granger

import (
	"context"
	"time"
)

func (g *Granger) TimelineRefresher(quit chan struct{}) {
	g.log.Dbg("TimelineRefresher started")

	timelineRefreshInterval := 1 * time.Minute

	for {
		select {
		case <-quit:
			return
		case <-time.Tick(timelineRefreshInterval):
			g.log.Inf("Refreshing timelines")
			err := g.db.RefreshOldestTimelines(context.Background())
			if err != nil {
				g.log.Err("Failed to refresh timelines", "error", err)
			}
		}
	}
}
