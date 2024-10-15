package granger

import (
	"context"
	"time"
)

func (g *Granger) cleanOldOnboardTokens(quit chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			g.log.Debug("cleanOldOnboardTokens quitting")
			return
		case <-time.After(1 * time.Minute):
			ctx := context.Background()
			_ = g.db.CleanOldOnboardTokens(ctx)
		}
	}
}
