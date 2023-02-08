package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/serialize"
	"electionguard-verifier-go/util"
	"fmt"
)

func main() {
	// Singleton
	cipherTextElectionRecord := serialize.ParseFromJsonToSingleObject(util.SAMPLE_DATA_DIR+"/context.json", schema.CiphertextElectionRecord{})
	manifest := serialize.ParseFromJsonToSingleObject(util.SAMPLE_DATA_DIR+"/manifest.json", schema.Manifest{})
	// cipertextTally := ...
	electionConstants := serialize.ParseFromJsonToSingleObject(util.SAMPLE_DATA_DIR+"/constants.json", schema.ElectionConstants{})
	// plaintextTally := ...

	// Non-singleton
	encryptionDevices := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR+"/encryption_devices/", schema.EncryptionDevice{})
	// spoiledBallots := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR + "/spoiled_ballots/")
	guardians := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR+"/guardians/", schema.Guardian{})

	// Verifying election data
	verifier := *core.MakeVerifier()
	electionIsValid := verifier.Verify(
		cipherTextElectionRecord,
		manifest,
		electionConstants,
		encryptionDevices,
		guardians)

	// Result
	if electionIsValid {
		fmt.Println("Election is valid")
	} else {
		fmt.Println("Election is invalid")
	}

}
