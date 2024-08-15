package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

func InitLogger() error {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err = config.Build()

	if err != nil {
		return err
	}

	defer logger.Sync()

	return nil
}

func GetLogger() *zap.Logger {
	return logger
}
