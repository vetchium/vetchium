package wand

import (
	"github.com/vetchium/vetchium/api/internal/config"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/postgres"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
)

// TODO: The Wand interface is needed because the alternative is to have every
// handler get the hermione pointer as a parameter, in addition to w and r,
// which is also not great. This Wand interface is also great. Should find a way
// to cleanly refactor this. Defining the handlers as methods on Hermione could
// be a step in that direction.
type Wand interface {
	DB() *postgres.PG
	Vator() *vetchi.Vator
	Hedwig() hedwig.Hedwig

	Config() *config.Hermione

	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Err(msg string, args ...any)
}
