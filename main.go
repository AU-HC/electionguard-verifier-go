package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/utility"
	"fmt"
)

func main() {
	// Fetching flags
	applicationArguments := utility.InitApplicationArguments()
	path := applicationArguments.ElectionArtifactsPath

	// Fetching logging level and creating logger
	loggingLevel := applicationArguments.LoggingLevel
	logger := utility.ConfigureLogger(loggingLevel)

	// Create verifier, and verify election data
	verifier := *core.MakeVerifier(logger)
	electionIsValid := verifier.Verify(path)

	// Result of verification of election data
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}
}
