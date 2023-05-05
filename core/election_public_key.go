package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateJointPublicKey(er *deserialize.ElectionRecord) {
	// Validate election public-key (Step 3) [ERROR IN SPEC SHEET FOR (3.B)]
	helper := MakeValidationHelper(v.logger, 3, "Election public-key validation")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	elgamalPublicKey := schema.MakeBigIntFromInt(1)
	for _, guardian := range er.Guardians {
		elgamalPublicKey = v.mulP(elgamalPublicKey, &guardian.ElectionPublicKey)
	}

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	computedExtendedBaseHash := crypto.HashElems(er.CiphertextElectionRecord.CryptoBaseHash, er.CiphertextElectionRecord.CommitmentHash)

	helper.addCheck(step3A, elgamalPublicKey.Compare(&er.CiphertextElectionRecord.ElgamalPublicKey))
	helper.addCheck(step3B, extendedBaseHash.Compare(computedExtendedBaseHash))

	v.helpers[helper.VerificationStep] = helper
}
