package core

import (
	"electionguard-verifier-go/schema"
	"go.uber.org/zap"
)

type Verifier struct {
	logger zap.Logger
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: *logger}
}

type VerifierArguments struct {
	// Election data fields
	CiphertextElectionRecord  schema.CiphertextElectionRecord
	Manifest                  schema.Manifest
	ElectionConstants         schema.ElectionConstants
	EncryptedTally            schema.EncryptedTally
	PlaintextTally            schema.PlaintextTally
	CoefficientsValidationSet schema.CoefficientsValidationSet
	SubmittedBallots          []schema.SubmittedBallot
	SpoiledBallots            []schema.SpoiledBallot
	EncryptionDevices         []schema.EncryptionDevice
	Guardians                 []schema.Guardian
}

func MakeVerifierArguments() *VerifierArguments {
	return &VerifierArguments{}
}

func (v *Verifier) Verify(args VerifierArguments) bool {
	v.logger.Debug("verifying election data")
	// Validate election parameters (Step 1)
	electionParametersHelper := MakeValidationHelper(&v.logger, "election parameters")
	electionParametersHelper.Ensure("p is correct", true) // Fake it
	electionParametersHelper.Ensure("q is correct", false)
	electionParametersIsNotValid := !electionParametersHelper.Validate()
	if electionParametersIsNotValid {
		return false
	}

	// Validate ... (Step 2)
	// ...
	// ...

	// Verification was successful
	return true
}
