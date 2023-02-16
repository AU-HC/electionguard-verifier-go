package core

import (
	"go.uber.org/zap"
)

type ValidationHelper struct {
	description string
	logger      *zap.Logger
	invariants  map[string]bool
}

func MakeValidationHelper(logger *zap.Logger, description string) *ValidationHelper {
	return &ValidationHelper{logger: logger, invariants: make(map[string]bool), description: description}
}

func (v *ValidationHelper) AddCheck(invariantDescription string, invariant bool) {
	v.invariants[invariantDescription] = invariant
}

func (v *ValidationHelper) Validate() bool {
	isValid := true
	errorMessages := make([]string, len(v.invariants))

	// Looping through each invariant and checking if it holds
	for description, invariant := range v.invariants {
		v.logger.Debug("Checked invariant: " + description)
		if !invariant {
			errorMessages = append(errorMessages, description)
			isValid = false
		}
	}

	if isValid {
		v.logger.Info("[VALID]: " + v.description)
	} else {
		v.logger.Info("[INVALID]: " + v.description)
		v.printAllErrors(errorMessages)
	}

	return isValid
}

func (v *ValidationHelper) printAllErrors(errors []string) {
	for _, err := range errors {
		if err != "" {
			v.logger.Info(err)
		}
	}
}
