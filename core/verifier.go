package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
	"go.uber.org/zap"
)

type Verifier struct {
	logger *zap.Logger
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger}
}

func (v *Verifier) Verify(path string) bool {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	args := parser.ConvertJsonDataToGoStruct(path)
	v.logger.Info("[VALID]: Loaded election data (Step 0)")

	// Validate election parameters (Step 1):
	correctConstants := utility.MakeCorrectElectionConstants()
	electionParametersHelper := MakeValidationHelper(v.logger, "Election parameters (Step 1)")
	electionParametersHelper.AddCheck("(1.A) The large prime is equal to the large modulus p", correctConstants.P.Compare(&args.ElectionConstants.LargePrime))
	electionParametersHelper.AddCheck("(1.B) The small prime is equal to the prime q", correctConstants.Q.Compare(&args.ElectionConstants.SmallPrime))
	electionParametersHelper.AddCheck("(1.C) The cofactor is equal to r = (p − 1)/q", correctConstants.C.Compare(&args.ElectionConstants.Cofactor))
	electionParametersHelper.AddCheck("(1.D) The generator is equal to the generator g", correctConstants.G.Compare(&args.ElectionConstants.Generator))
	electionParametersIsNotValid := !electionParametersHelper.Validate()
	if electionParametersIsNotValid {
		return false
	}

	// Validate guardian public-key (Step 2)
	publicKeyValidationHelper := MakeValidationHelper(v.logger, "Guardian public-key validation (Step 2)")
	/* for _, guardian := range args.Guardians {
		for _, _ := range guardian.ElectionCommitments {

		}
	}*/
	publicKeysIsNotValid := !publicKeyValidationHelper.Validate()
	if publicKeysIsNotValid {
		return false
	}

	// Verification was successful
	return true
}
