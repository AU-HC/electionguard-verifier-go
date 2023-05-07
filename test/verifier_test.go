package test

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/utility"
	"runtime"
	"testing"
)

func getSut() core.Verifier {
	// Setup code
	// Fetching flags and amount of logical cores for the CPU
	amountOfLogicalCores := runtime.NumCPU()

	// Fetching logging level and creating logger
	logger := utility.ConfigureLogger(0)
	outputStrategy := core.NoOutputStrategy{}
	verifyStrategy := core.MakeVerifyStrategy(true, amountOfLogicalCores)

	// Create verifier, set strategies, and verify election data
	verifier := *core.MakeVerifier(logger)
	verifier.SetOutputStrategy(outputStrategy)
	verifier.SetVerifyStrategy(verifyStrategy)

	return verifier
}

func TestDataWillVerify(t *testing.T) {
	// Setup
	sut := getSut()
	path := "./data/full_valid/election_record"

	// Act
	electionIsValid := sut.Verify(path)

	// Assert
	if !electionIsValid {
		t.Error("Expected result to be valid, but was invalid")
	}
}

func TestDataWillNotVerify(t *testing.T) {
	// Setup
	sut := getSut()
	path := "./data/full_invalid/election_record"

	// Act
	electionIsValid := sut.Verify(path)

	// Assert
	if electionIsValid {
		t.Error("Expected result to be invalid, but was valid")
	}
}
