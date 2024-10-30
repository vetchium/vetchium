package vhandler

import (
	"log/slog"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// Might as well call it a Wand
type VHandler interface {
	DB() db.DB
	Log() *slog.Logger
	Vator() *vetchi.Vator

	// TODO: Implement a better way to pass config to the handlers
	TGTLife() time.Duration
}
