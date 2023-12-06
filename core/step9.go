package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateTallyDecryptions(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 9, "Correctness of tally decryptions")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	g := er.ElectionConstants.Generator
	q := &er.ElectionConstants.SmallPrime
	ctx := er.CiphertextElectionRecord
	k := ctx.ElgamalPublicKey
	extendedBaseHash := ctx.CryptoExtendedBaseHash

	for _, contest := range er.PlaintextTally.Contests {
		encryptedContest := er.EncryptedTally.Contests[contest.ObjectId]

		for _, selection := range contest.Selections {
			encryptedSelection := encryptedContest.Selections[selection.ObjectId]

			errorString := "(ContestID:" + contest.ObjectId + ", SelectionID:" + selection.ObjectId + ")"
			helper.addCheck("(9.A) The response is not valid.", v.isInRange(selection.Proof.Response), errorString)

			// Computing values needed for 9.C
			m := v.mulP(&encryptedSelection.Ciphertext.Data, v.invP(&selection.Value))
			a := v.mulP(v.powP(&g, &selection.Proof.Response), v.powP(&k, &selection.Proof.Challenge))
			b := v.mulP(v.powP(&encryptedSelection.Ciphertext.Pad, &selection.Proof.Response), v.powP(m, &selection.Proof.Challenge))
			hash := crypto.Hash(q, "30")
			hash = crypto.Hash(q, hash, extendedBaseHash)
			hash = crypto.Hash(q, hash, k)
			hash = crypto.Hash(q, hash, encryptedSelection.Ciphertext.Pad)
			hash = crypto.Hash(q, hash, encryptedSelection.Ciphertext.Data)
			hash = crypto.Hash(q, hash, a)
			hash = crypto.Hash(q, hash, b)
			hash = crypto.Hash(q, hash, m)

			helper.addCheck("(9.C) The challenge is not computed correctly.", selection.Proof.Challenge.Compare(hash), errorString)
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
