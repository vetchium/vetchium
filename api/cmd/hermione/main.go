package main

import (
	"log/slog"
	"os"

	"github.com/psankar/vetchi/api/internal/hermione"
)

func main() {
	slog.Info("Hermione starting up ...")

	hermione, err := hermione.NewHermione()
	if err != nil {
		slog.Error("Failed to create new Hermione", "error", err)
		os.Exit(1)
	}

	if err := hermione.Run(); err != nil {
		slog.Error("Failed to run Hermione", "error", err)
		os.Exit(1)
	}
}
