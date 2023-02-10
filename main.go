package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/serialize"
	"electionguard-verifier-go/utility"
	"fmt"
)

func main() {
	// Logger
	logger := utility.ConfigureLogger()

	// Singleton files
	cipherTextElectionRecord := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/context.json", schema.CiphertextElectionRecord{})
	manifest := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/manifest.json", schema.Manifest{})
	encryptedTally := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/encrypted_tally.json", schema.EncryptedTally{})
	electionConstants := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/constants.json", schema.ElectionConstants{})
	plaintextTally := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/tally.json", schema.PlaintextTally{})
	coefficients := serialize.ParseFromJsonToSingleObject(utility.SAMPLE_DATA_DIR+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	encryptionDevices := serialize.ParseFromJsonToSlice(utility.SAMPLE_DATA_DIR+"/encryption_devices/", schema.EncryptionDevice{})
	spoiledBallots := serialize.ParseFromJsonToSlice(utility.SAMPLE_DATA_DIR+"/spoiled_ballots/", schema.SpoiledBallot{})
	submittedBallots := serialize.ParseFromJsonToSlice(utility.SAMPLE_DATA_DIR+"/submitted_ballots/", schema.SubmittedBallots{})
	guardians := serialize.ParseFromJsonToSlice(utility.SAMPLE_DATA_DIR+"/guardians/", schema.Guardian{})

	// Create arguments for verifier
	verifierArguments := *core.MakeVerifierArguments()
	verifierArguments.CiphertextElectionRecord = cipherTextElectionRecord
	verifierArguments.Manifest = manifest
	verifierArguments.EncryptedTally = encryptedTally
	verifierArguments.ElectionConstants = electionConstants
	verifierArguments.PlaintextTally = plaintextTally
	verifierArguments.CoefficientsValidationSet = coefficients
	verifierArguments.EncryptionDevices = encryptionDevices
	verifierArguments.SpoiledBallots = spoiledBallots
	verifierArguments.SubmittedBallots = submittedBallots
	verifierArguments.Guardians = guardians
	verifierArguments.Logger = logger

	// Creating verifier and verifying election data
	verifier := *core.MakeVerifier()
	electionIsValid := verifier.Verify(verifierArguments)

	// Result of verification of election data
	fmt.Println()
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}
}
