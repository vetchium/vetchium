package vhandler

import (
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// Might as well call it a Wand
type VHandler interface {
	DB() db.DB
	Vator() *vetchi.Vator

	Err(msg string, args ...any)
	Dbg(msg string, args ...any)
	Inf(msg string, args ...any)

	// TODO: Implement a better way to pass config to the handlers
	TGTLife() time.Duration
}
