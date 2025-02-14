package granger

import (
	"context"
	"time"
)

func (g *Granger) pruneOfficialEmailCodes(quit chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			return
		case <-time.Tick(5 * time.Minute):
			// TODO: Read the time interval from config
			err := g.db.PruneOfficialEmailCodes(context.Background())
			if err != nil {
				g.log.Err("failed to prune official email codes", "error", err)
			}
		}
	}
}
