package main

import (
	"log/slog"
	"os"

	"github.com/rhysmcneill/dividr/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("Fatal application error", "error", err)
		os.Exit(1)
	}
}
