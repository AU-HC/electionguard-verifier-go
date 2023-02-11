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
		return zap.NewProduction()
	case LogNone:
		return zap.NewNop(), nil
	}

	return nil, fmt.Errorf("no valid logging level provided")
}
