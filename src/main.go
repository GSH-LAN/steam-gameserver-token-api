package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/gsh-lan/steam-gameserver-token-api/src/app"
	"github.com/gsh-lan/steam-gameserver-token-api/src/environment"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		// set level from environment "LOG_LEVEL"
		Level: func() (lvl slog.Level) {
			lvl = slog.LevelInfo
			if err := lvl.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
				slog.Error("Invalid LOG_LEVEL, defaulting to INFO", "error", err)
			}
			return
		}(),
	}))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		slog.Error("Error loading .env file", "error", err)
		panic(err)
	}

	e := environment.Load()

	a := app.App{}

	a.Run(e.SteamWebAPIBindAddress, e.SteamWebAPIKey, e.AuthToken, e.BackgroundProcessingInterval)
}
