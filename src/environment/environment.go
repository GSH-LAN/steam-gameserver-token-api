package environment

import (
	"github.com/caarlos0/env/v10"
	"github.com/gsh-lan/steam-gameserver-token-api/src/logger"
	"go.uber.org/zap"
	"time"
)

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

type Environment struct {
	SteamWebAPIBindAddress       string        `env:"STEAM_WEB_API_BIND_ADDRESS" envDefault:":8080"`
	SteamWebAPIKey               string        `env:"STEAM_WEB_API_KEY"`
	AuthToken                    string        `env:"AUTH_TOKEN"`
	BackgroundProcessingInterval time.Duration `env:"BACKGROUND_PROCESSING_INTERVAL" envDefault:"10s"`
}

func Load() *Environment {
	e := Environment{}
	if err := env.Parse(&e); err != nil {
		log.Panicf("%+v\n", err)
	}

	log.Infof("Initialized Environment: %+v", e)
	return &e
}
