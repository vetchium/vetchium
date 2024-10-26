package vhandler

import (
	"log/slog"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// Might as well call it a Wand
type VHandler interface {
	DB() db.DB
	Log() *slog.Logger
	Vator() *vetchi.Vator
}
