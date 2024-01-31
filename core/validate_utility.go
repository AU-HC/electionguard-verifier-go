package core

import (
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type ValidationHelper struct {
	VerificationStep int
	Description      string
	Checked, Failed  int
	TimeToVerifyInMs int64
	ErrorMessage     []string
	isValid          bool
	logger           *zap.Logger
	mu               sync.Mutex
	wg               sync.WaitGroup
}

func MakeValidationHelper(logger *zap.Logger, step int, description string) *ValidationHelper {
	return &ValidationHelper{
		logger:           logger,
		Description:      description,
		VerificationStep: step,
		isValid:          true,
		mu:               sync.Mutex{},
		wg:               sync.WaitGroup{},
	}
}

func (v *ValidationHelper) addCheck(invariantDescription string, invariant bool, errorString ...string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// If invariant is true, do nothing
	v.logger.Debug("Checked invariant: " + invariantDescription)
	v.Checked += 1
	if invariant {
		return
	}

	// else append the error message and increment failed invariants
	v.isValid = false
	v.Failed += 1
	errorMessage := invariantDescription
	for _, s := range errorString { // only zero or one string is supplied (done as an optional argument)
		errorMessage += s
	}
	v.ErrorMessage = append(v.ErrorMessage, errorMessage)

}

func (v *ValidationHelper) validate() bool {
	stepString := strconv.Itoa(v.VerificationStep)

	if v.isValid {
		v.logger.Info("[VALID]: " + stepString + ". " + v.Description)
		return true
	}

	v.logger.Info("[INVALID]: " + stepString + ". " + v.Description)

	for _, errorString := range v.ErrorMessage {
		v.logger.Debug(errorString)
	}
	return false
}

func (v *ValidationHelper) measureTimeToValidateStep(start time.Time) {
	total := time.Since(start)
	v.TimeToVerifyInMs = total.Milliseconds()
	v.logger.Info("Validation of step " + strconv.Itoa(v.VerificationStep) + " took: " + total.String())
}
