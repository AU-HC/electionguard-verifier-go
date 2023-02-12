package core

import (
	"bytes"
	"go.uber.org/zap"
)

type ValidationHelper struct {
	description string
	logger      zap.Logger
	invariants  map[string]bool
}

func MakeValidationHelper(logger *zap.Logger, description string) *ValidationHelper {
	return &ValidationHelper{logger: *logger, invariants: make(map[string]bool), description: description}
}

func (v *ValidationHelper) Ensure(invariantDescription string, invariant bool) {
	v.invariants[invariantDescription] = invariant
}

func (v *ValidationHelper) Validate() bool {
	isValid := true
	var errorMsg bytes.Buffer

	// Looping through each invariant and checking if it holds
	for description, invariant := range v.invariants {
		if !invariant {
			errorMsg.WriteString("failed to validate invariant: " + description + "\n")
			isValid = false
		}
	}

	if isValid {
		v.logger.Info("valid: " + v.description)
	} else {
		v.logger.Info("invalid: " + v.description)
		v.logger.Info(errorMsg.String())
	}

	return isValid
}
