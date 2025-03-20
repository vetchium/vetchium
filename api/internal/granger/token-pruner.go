package granger

import (
	"context"
	"time"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (g *Granger) pruneTokens(quit chan struct{}) {
	g.log.Dbg("Starting pruneTokens job")
	defer g.log.Dbg("pruneTokens job finished")
	defer g.wg.Done()

	ticker := time.NewTicker(vetchi.PruneTokensInterval)

	for {
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("pruneTokens quitting")
			return
		case <-ticker.C:
			ticker.Stop()
			_ = g.db.PruneTokens(context.Background())
			ticker = time.NewTicker(vetchi.PruneTokensInterval)
		}
	}
}
