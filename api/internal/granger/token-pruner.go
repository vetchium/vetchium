package granger

import (
	"context"
	"time"
)

func (g *Granger) pruneTokens(quit chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			g.log.Debug("pruneTokens quitting")
			return
		case <-time.After(1 * time.Minute):
			_ = g.db.PruneTokens(context.Background())
		}
	}
}
