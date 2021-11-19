package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	var config = zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	l, err := config.Build()
	if err != nil {
		fmt.Printf("Error happened initializing logger: %+v", err)
		panic(err)
	}
	logger = l.Sugar()
}

func GetSugaredLogger() *zap.SugaredLogger {
	return logger
}
