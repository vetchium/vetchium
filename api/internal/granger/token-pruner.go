package granger

import (
	"context"
	"time"
)

func (g *Granger) pruneTokens(quit chan struct{}) {
	g.log.Dbg("Starting pruneTokens job")
	defer g.log.Dbg("pruneTokens job finished")
	defer g.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)

	for {
		ticker.Reset(1 * time.Minute)
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("pruneTokens quitting")
			return
		case <-ticker.C:
			ticker.Stop()
			_ = g.db.PruneTokens(context.Background())
		}
	}
}
