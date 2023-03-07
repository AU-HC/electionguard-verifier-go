package core

import (
	"go.uber.org/zap"
	"strings"
)

type ValidationHelper struct {
	description, step string
	logger            *zap.Logger
	errorMsg          *strings.Builder
	checked, failed   int
}

func MakeValidationHelper(logger *zap.Logger, description string) *ValidationHelper {
	return &ValidationHelper{logger: logger, description: description, errorMsg: &strings.Builder{}}
}

func (v *ValidationHelper) addCheck(invariantDescription string, invariant bool) {
	// If invariant is true, do nothing
	v.logger.Debug("Checked invariant: " + invariantDescription)
	v.checked += 1
	if invariant {
		return
	}

	// else append the error message and increment failed invariants
	v.failed += 1
	v.errorMsg.WriteString(invariantDescription)
	v.errorMsg.WriteString("\n")
}

func (v *ValidationHelper) validate() bool {
	if v.errorMsg.Len() == 0 {
		v.logger.Info("[VALID]: " + v.description)
		return true
	}

	v.logger.Info("[INVALID]: " + v.description)
	v.logger.Debug(v.errorMsg.String())
	return false
}
