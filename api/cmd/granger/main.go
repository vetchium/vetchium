package main

import (
	"log/slog"
	"os"

	"vetchi.org/internal/granger"
)

func main() {
	slog.Info("Granger starting up ...")

	granger, err := granger.NewGranger()
	if err != nil {
		slog.Error("Failed to initialize Granger", "error", err)
		os.Exit(1)
	}

	if err := granger.Run(); err != nil {
		slog.Error("Failed to run Granger", "error", err)
		os.Exit(1)
	}
}
