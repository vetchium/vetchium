package main

import (
	"log/slog"
	"os"

	"vetchi.org/internal/hermione"
)

func main() {
	slog.Info("Hermione launching ...")

	hermione, err := hermione.New()
	if err != nil {
		slog.Error("Failed to create new Hermione", "error", err)
		os.Exit(1)
	}

	if err := hermione.Run(); err != nil {
		slog.Error("Failed to run Hermione", "error", err)
		os.Exit(1)
	}
}
