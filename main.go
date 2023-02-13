package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/deserialize"
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

	// Create verifier, parser and arguments for verifier
	verifier := *core.MakeVerifier(logger)
	parser := *deserialize.MakeParser(logger)
	electionData := parser.ConvertJsonDataToGoStruct(path)

	// Verifying election data
	electionIsValid := verifier.Verify(electionData)

	// Result of verification of election data
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}
}
