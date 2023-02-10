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

type VerifierArguments struct {
	CipherTextElectionRecord  schema.CiphertextElectionRecord
	Manifest                  schema.Manifest
	ElectionConstants         schema.ElectionConstants
	EncryptedTally            schema.EncryptedTally
	PlaintextTally            schema.PlaintextTally
	CoefficientsValidationSet schema.CoefficientsValidationSet
	SubmittedBallots          []schema.SubmittedBallots
	SpoiledBallots            []schema.SpoiledBallot
	EncryptionDevices         []schema.EncryptionDevice
	Guardians                 []schema.Guardian
}

func MakeVerifierArguments() *VerifierArguments {
	return &VerifierArguments{}
}

func (v *Verifier) Verify(args VerifierArguments) bool {
	fmt.Println(args.CoefficientsValidationSet.Coefficients["4"])

	return true
}
