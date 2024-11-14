package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Initialize logger
func InitLogger() *zap.Logger {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	productionConfig := zap.NewProductionConfig()
	productionConfig.Encoding = "console"
	productionConfig.Level = level
	productionConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := productionConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
