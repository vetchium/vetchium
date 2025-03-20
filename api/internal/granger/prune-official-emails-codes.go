package granger

import (
	"context"
	"time"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (g *Granger) pruneOfficialEmailCodes(quit chan struct{}) {
	g.log.Dbg("Starting pruneOfficialEmailCodes job")
	defer g.log.Dbg("pruneOfficialEmailCodes job finished")
	defer g.wg.Done()

	ticker := time.NewTicker(vetchi.PruneOfficialEmailCodesInterval)

	for {
		select {
		case <-quit:
			ticker.Stop()
			g.log.Dbg("pruneOfficialEmailCodes quitting")
			return
		case <-ticker.C:
			ticker.Stop()
			err := g.db.PruneOfficialEmailCodes(context.Background())
			if err != nil {
				g.log.Err("failed to prune official email codes", "error", err)
			}
			ticker = time.NewTicker(vetchi.PruneOfficialEmailCodesInterval)
		}
	}
}
