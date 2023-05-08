package core

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ValidationHelper struct {
	VerificationStep int
	Description      string
	Checked, Failed  int
	TimeToVerify     int64 // in ms
	ErrorMessage     string
	isValid          bool
	logger           *zap.Logger
	errorMsg         *strings.Builder
	mu               sync.Mutex
	wg               sync.WaitGroup
}

func MakeValidationHelper(logger *zap.Logger, step int, description string) *ValidationHelper {
	return &ValidationHelper{
		logger:           logger,
		Description:      description,
		errorMsg:         &strings.Builder{},
		VerificationStep: step,
		isValid:          true,
		mu:               sync.Mutex{},
		wg:               sync.WaitGroup{},
	}
}

func (v *ValidationHelper) addCheck(invariantDescription string, invariant bool) {
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
	v.errorMsg.WriteString(invariantDescription)
	v.errorMsg.WriteString("\n")
}

func (v *ValidationHelper) validate() bool {
	stepString := strconv.Itoa(v.VerificationStep)

	if v.isValid {
		v.logger.Info("[VALID]: " + stepString + ". " + v.Description)
		return true
	}

	v.logger.Info("[INVALID]: " + stepString + ". " + v.Description)
	v.logger.Debug(v.errorMsg.String())

	v.ErrorMessage = v.errorMsg.String()
	return false
}

func (v *ValidationHelper) measureTimeToValidateStep(start time.Time) {
	total := time.Since(start)
	v.TimeToVerify = total.Milliseconds()
	v.logger.Info("Validation of step " + strconv.Itoa(v.VerificationStep) + " took: " + total.String())
}
