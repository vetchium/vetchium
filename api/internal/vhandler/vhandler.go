package vhandler

import (
	"time"

	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/postgres"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// Might as well call it a Wand
type VHandler interface {
	DB() *postgres.PG
	Vator() *vetchi.Vator

	// TODO: Implement a better way to pass config to the handlers
	TGTLife() time.Duration

	Hedwig() hedwig.Hedwig

	util.Logger
}
