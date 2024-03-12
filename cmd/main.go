package main

import (
	"bot21/internal/app"
	"log/slog"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error(err.Error())
	}
}
