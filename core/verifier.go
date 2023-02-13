package core

import (
	"electionguard-verifier-go/schema"
	"go.uber.org/zap"
)

type Verifier struct {
	logger *zap.Logger
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger}
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

	// Validate election parameters (Step 1):
	correctConstants := MakeCorrectElectionConstants()
	electionParametersHelper := MakeValidationHelper(v.logger, "Election parameters (Step 1)")
	electionParametersHelper.AddCheck("(1.A) The large prime is equal to the large modulus p",
		correctConstants.P.Compare(&args.ElectionConstants.LargePrime))
	electionParametersHelper.AddCheck("(1.B) The small prime is equal to the prime q",
		correctConstants.Q.Compare(&args.ElectionConstants.SmallPrime))
	electionParametersHelper.AddCheck("(1.C) The cofactor is equal to r = (p − 1)/q",
		correctConstants.C.Compare(&args.ElectionConstants.Cofactor))
	electionParametersHelper.AddCheck("(1.D) The generator is equal to the generator g",
		correctConstants.G.Compare(&args.ElectionConstants.Generator))
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
