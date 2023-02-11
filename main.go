package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/serialize"
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

	// Create verifier and arguments for verifier
	verifier := *core.MakeVerifier(logger)
	verifierArguments := *core.MakeVerifierArguments()

	// Parsing singleton files
	verifierArguments.CiphertextElectionRecord = serialize.ParseFromJsonToSingleObject(path+"/context.json", schema.CiphertextElectionRecord{})
	verifierArguments.Manifest = serialize.ParseFromJsonToSingleObject(path+"/manifest.json", schema.Manifest{})
	verifierArguments.EncryptedTally = serialize.ParseFromJsonToSingleObject(path+"/encrypted_tally.json", schema.EncryptedTally{})
	verifierArguments.ElectionConstants = serialize.ParseFromJsonToSingleObject(path+"/constants.json", schema.ElectionConstants{})
	verifierArguments.PlaintextTally = serialize.ParseFromJsonToSingleObject(path+"/tally.json", schema.PlaintextTally{})
	verifierArguments.CoefficientsValidationSet = serialize.ParseFromJsonToSingleObject(path+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	verifierArguments.EncryptionDevices = serialize.ParseFromJsonToSlice(path+"/encryption_devices/", schema.EncryptionDevice{})
	verifierArguments.SpoiledBallots = serialize.ParseFromJsonToSlice(path+"/spoiled_ballots/", schema.SpoiledBallot{})
	verifierArguments.SubmittedBallots = serialize.ParseFromJsonToSlice(path+"/submitted_ballots/", schema.SubmittedBallots{})
	verifierArguments.Guardians = serialize.ParseFromJsonToSlice(path+"/guardians/", schema.Guardian{})

	// Verifying election data
	electionIsValid := verifier.Verify(verifierArguments)

	// Result of verification of election data
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}
}
