package main

import (
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/serialize"
	"fmt"
)

func main() {
	// Parse election data
	cipherTextElectionRecord := serialize.ParseFromJson("data/hamilton-general/election_record/context.json", schema.CiphertextElectionRecord{})
	manifest := serialize.ParseFromJson("data/hamilton-general/election_record/manifest.json", schema.Manifest{})

	// Test print
	fmt.Println(cipherTextElectionRecord.ElgamalPublicKey)
	fmt.Println(manifest.Type)
}
