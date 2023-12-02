package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateExtendedBaseHash(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 4, "Extended base hash is correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	q := er.ElectionConstants.SmallPrime
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	baseHash := er.CiphertextElectionRecord.CryptoBaseHash
	commitmentHash := er.CiphertextElectionRecord.CommitmentHash
	electionPublicKey := er.CiphertextElectionRecord.ElgamalPublicKey

	// Calculating hash and comparing it
	calculatedExtendedBaseHash := crypto.Hash1(&q, baseHash, "12", electionPublicKey, commitmentHash)
	helper.addCheck("(4.A) The extended base hash is not computed correctly.", calculatedExtendedBaseHash.Compare(&extendedBaseHash))

	v.helpers[helper.VerificationStep] = helper
}
