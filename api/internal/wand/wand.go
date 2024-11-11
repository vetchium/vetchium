package wand

import (
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type Wand interface {
	DB() *postgres.PG
	Vator() *vetchi.Vator
	Hedwig() hedwig.Hedwig

	ConfigDuration(key db.TokenType) (time.Duration, error)

	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)
	Err(msg string, args ...any)
}
