package logging

import (
	"electionguard-verifier-go/error_handling"
	"go.uber.org/zap"
)

func ConfigureLogger(level Level) *zap.Logger {
	// Creating logger and checking for error
	logger, err := createLogger(level)
	error_handling.PanicError(err)

	// Created logger and returning
	logger.Debug("successfully created logger")

	return logger
}

func createLogger(level Level) (*zap.Logger, error) {
	switch level {
	case LogDebug:
		return zap.NewDevelopment()
	case LogInfo: // Changed Config to have same style as development
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
		return config.Build()
	default:
		return zap.NewNop(), nil
	}
}
