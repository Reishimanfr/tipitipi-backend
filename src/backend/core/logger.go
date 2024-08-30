package core

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	err    error
)

func InitLogger() (*zap.Logger, error) {
	var config zap.Config

	if os.Getenv("DEV") != "true" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	Logger, err = config.Build()

	if err != nil {
		return nil, err
	}

	defer Logger.Sync()

	return Logger, nil
}

func GetLogger() *zap.Logger {
	return Logger
}
