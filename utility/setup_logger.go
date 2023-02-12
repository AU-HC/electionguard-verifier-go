package utility

import (
	"fmt"
	"go.uber.org/zap"
)

func ConfigureLogger(level LoggingLevel) *zap.Logger {
	// Creating logger and checking for error
	logger, err := createLogger(level)
	PanicError(err)

	// Created logger and returning
	logger.Debug("successfully created logger")

	return logger
}

func createLogger(level LoggingLevel) (*zap.Logger, error) {
	switch level {
	case LogDebug:
		return zap.NewDevelopment()
	case LogInfo:
		xd := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "console",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		return xd.Build()
	case LogNone:
		return zap.NewNop(), nil
	}

	return nil, fmt.Errorf("no valid logging level provided")
}
