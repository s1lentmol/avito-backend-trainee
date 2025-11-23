package main

import (
	"log/slog"
	"os"

	"github.com/silentmol/avito-backend-trainee/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app listen", slog.Any("error", err))
		os.Exit(1)
	}
}