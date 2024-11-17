package wand

import (
	"github.com/psankar/vetchi/api/internal/config"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/vetchi"
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
