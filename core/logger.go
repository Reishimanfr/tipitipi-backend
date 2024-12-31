package core

import (
	"bash06/tipitipi-backend/flags"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() (*zap.Logger, error) {
	var config zap.Config

	if *flags.Dev {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()

	if err != nil {
		return nil, err
	}

	defer logger.Sync()

	return logger, nil
}
