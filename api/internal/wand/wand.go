package wand

import (
	"github.com/vetchium/vetchium/api/internal/config"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/postgres"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
)

type Wand interface {
	DB() *postgres.PG
	Vator() *vetchi.Vator
	Hedwig() hedwig.Hedwig

	Config() *config.Hermione

	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Err(msg string, args ...any)
}
