package core

import (
	"electionguard-verifier-go/schema"
	"go.uber.org/zap"
)

type Verifier struct {
}

func MakeVerifier() *Verifier {
	return &Verifier{}
}

type VerifierArguments struct {
	// Logger
	Logger *zap.Logger
	// Election data fields
	CiphertextElectionRecord  schema.CiphertextElectionRecord
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
	// Fetch logger from arguments
	logger := *args.Logger
	logger.Info("log from verifier")

	// Validate election parameters

	// Verification was successful
	return true
}
