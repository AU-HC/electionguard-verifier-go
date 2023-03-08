package core

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type ValidationHelper struct {
	description                       string
	logger                            *zap.Logger
	errorMsg                          *strings.Builder
	checked, failed, verificationStep int
}

func MakeValidationHelper(logger *zap.Logger, step int, description string) *ValidationHelper {
	return &ValidationHelper{logger: logger, description: description, errorMsg: &strings.Builder{}, verificationStep: step}
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
	stepString := strconv.Itoa(v.verificationStep)

	if v.errorMsg.Len() == 0 {
		v.logger.Info("[VALID]: " + stepString + ". " + v.description)
		return true
	}

	v.logger.Info("[INVALID]: " + stepString + ". " + v.description)
	v.logger.Debug(v.errorMsg.String())
	return false
}
