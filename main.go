package main

import (
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
	// plaintextTally :=

	// Non-singleton
	encryptionDevices := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR+"/encryption_devices/", schema.EncryptionDevice{})
	// spoiledBallots := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR + "/spoiled_ballots/")
	guardians := serialize.ParseFromJsonToSlice(util.SAMPLE_DATA_DIR+"/guardians/", schema.Guardian{})

	// Test print
	fmt.Println(cipherTextElectionRecord.Configuration.MaxVotes)
	fmt.Println(manifest.Type)
	fmt.Println(electionConstants.SmallPrime)
	fmt.Println(encryptionDevices[0])
	fmt.Print(guardians[0].ElectionProofs[0].Commitment)
}
