package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/serialize"
	"electionguard-verifier-go/utility"
	"fmt"
)

func main() {
	// Creating logger
	logger := utility.ConfigureLogger(utility.LogDebug)

	// Create arguments for verifier
	verifierArguments := *core.MakeVerifierArguments()

	// Singleton files
	verifierArguments.CiphertextElectionRecord = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/context.json", schema.CiphertextElectionRecord{})
	verifierArguments.Manifest = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/manifest.json", schema.Manifest{})
	verifierArguments.EncryptedTally = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/encrypted_tally.json", schema.EncryptedTally{})
	verifierArguments.ElectionConstants = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/constants.json", schema.ElectionConstants{})
	verifierArguments.PlaintextTally = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/tally.json", schema.PlaintextTally{})
	verifierArguments.CoefficientsValidationSet = serialize.ParseFromJsonToSingleObject(utility.SampleDataDir+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	verifierArguments.EncryptionDevices = serialize.ParseFromJsonToSlice(utility.SampleDataDir+"/encryption_devices/", schema.EncryptionDevice{})
	verifierArguments.SpoiledBallots = serialize.ParseFromJsonToSlice(utility.SampleDataDir+"/spoiled_ballots/", schema.SpoiledBallot{})
	verifierArguments.SubmittedBallots = serialize.ParseFromJsonToSlice(utility.SampleDataDir+"/submitted_ballots/", schema.SubmittedBallots{})
	verifierArguments.Guardians = serialize.ParseFromJsonToSlice(utility.SampleDataDir+"/guardians/", schema.Guardian{})

	// Adding logger to arguments
	verifierArguments.Logger = logger

	// Creating verifier and verifying election data
	verifier := *core.MakeVerifier()
	electionIsValid := verifier.Verify(verifierArguments)

	// Result of verification of election data
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}
}
