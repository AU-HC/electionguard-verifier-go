package core

import (
	"electionguard-verifier-go/schema"
)

type Verifier struct {
}

func MakeVerifier() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(cipherTextElectionRecord schema.CiphertextElectionRecord, manifest schema.Manifest, electionConstants schema.ElectionConstants, encryptionDevices []schema.EncryptionDevice, guardians []schema.Guardian, encryptedTally schema.EncryptedTally, submittedBallots []schema.SubmittedBallots, spoiledBallots []schema.SpoiledBallot) bool {
	return true
}
