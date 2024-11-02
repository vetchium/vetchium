package wand

import (
	"time"

	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Wand interface {
	DB() *postgres.PG
	Vator() *vetchi.Vator
	Hedwig() hedwig.Hedwig

	// TODO: Implement a better way to pass config to the handlers
	TGTLife() time.Duration

	// Embed Err Dbg Inf methods instead of Log().Err() so that code is narrower
	util.Logger
}
