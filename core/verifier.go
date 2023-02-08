package core

import (
	"electionguard-verifier-go/schema"
	"fmt"
)

type Verifier struct {
}

func MakeVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(cipherTextElectionRecord schema.CiphertextElectionRecord,
	manifest schema.Manifest,
	electionConstants schema.ElectionConstants,
	encryptionDevices []schema.EncryptionDevice,
	guardians []schema.Guardian) bool {

	// Test print
	fmt.Println(cipherTextElectionRecord.Configuration.MaxVotes)
	fmt.Println(manifest.Type)
	fmt.Println(electionConstants.SmallPrime)
	fmt.Println(encryptionDevices[0])
	// fmt.Print(guardians[0].ElectionProofs[0].Commitment)

	return true
}
