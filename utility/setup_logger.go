package utility

import (
	"go.uber.org/zap"
)

func ConfigureLogger() *zap.Logger {
	// Creating logger
	logger, err := zap.NewDevelopment()

	// Checking for error
	if err != nil {
		panic(err)
	}

	// Created logger and returning
	logger.Debug("successfully created logger")

	return logger
}
