package postgres

import (
	"context"

	"github.com/vetchium/vetchium/typespec/hub"
)

func (pg *PG) GetMyHomeTimeline(
	ctx context.Context,
	req hub.GetMyHomeTimelineRequest,
) (hub.MyHomeTimeline, error) {
	pg.log.Dbg("Entered PG GetMyHomeTimeline")

	return hub.MyHomeTimeline{}, nil
}
