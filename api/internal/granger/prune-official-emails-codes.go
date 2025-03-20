package granger

import (
	"context"
	"time"
)

func (g *Granger) pruneOfficialEmailCodes(quit chan struct{}) {
	g.log.Dbg("Starting pruneOfficialEmailCodes job")
	defer g.log.Dbg("pruneOfficialEmailCodes job finished")
	defer g.wg.Done()

	ticker := time.NewTicker(5 * time.Minute)

	for {
		ticker.Reset(5 * time.Minute)
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("pruneOfficialEmailCodes quitting")
			return
		case <-ticker.C:
			ticker.Stop()
			g.log.Dbg("pruneOfficialEmailCodes ticker received signal")
			// TODO: Read the time interval from config
			err := g.db.PruneOfficialEmailCodes(context.Background())
			if err != nil {
				g.log.Err("failed to prune official email codes", "error", err)
			}
		}
	}
}
