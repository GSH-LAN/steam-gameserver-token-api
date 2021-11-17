package main

import (
	"strings"

	"github.com/gsh-lan/steam-gameserver-token-api/src/app"
	"github.com/gsh-lan/steam-gameserver-token-api/src/environment"
	"github.com/gsh-lan/steam-gameserver-token-api/src/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

func main() {
	err := godotenv.Load()
	if err != nil && !strings.Contains(err.Error(), "no such file") {
		log.Fatal("Error loading .env file: %+v", err)
	}

	e := environment.Load()

	a := app.App{}

	a.Run(e.SteamWebAPIBindAddress, e.SteamWebAPIKey, e.AuthToken)
}
