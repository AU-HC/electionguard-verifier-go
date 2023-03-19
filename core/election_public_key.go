package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateJointPublicKey(er *deserialize.ElectionRecord) {
	// Validate election public-key (Step 3) [ERROR IN SPEC SHEET FOR (3.B)]
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 3, "Election public-key validation")
	start := time.Now()

	elgamalPublicKey := schema.MakeBigIntFromString("1", 10)
	for _, guardian := range er.Guardians {
		elgamalPublicKey = v.mulP(elgamalPublicKey, &guardian.ElectionPublicKey)
	}

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	computedExtendedBaseHash := crypto.HashElems(er.CiphertextElectionRecord.CryptoBaseHash, er.CiphertextElectionRecord.CommitmentHash)

	helper.addCheck("(3.A) The joint public election key is computed correctly", elgamalPublicKey.Compare(&er.CiphertextElectionRecord.ElgamalPublicKey))
	helper.addCheck("(3.B) The extended base hash is computed correctly", extendedBaseHash.Compare(computedExtendedBaseHash))

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 3 took: " + time.Since(start).String())
}
