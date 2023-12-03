package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateTallyDecryptions(er *deserialize.ElectionRecord) {
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

			helper.addCheck("(9.A) The challenge is not valid.", v.isInRange(selection.Proof.Response))

			// Computing values needed for 9.C
			m := v.mulP(&encryptedSelection.Ciphertext.Data, v.invP(&selection.Value))
			a := v.mulP(v.powP(&g, &selection.Proof.Response), v.powP(&k, &selection.Proof.Challenge))
			b := v.mulP(v.powP(&encryptedSelection.Ciphertext.Pad, &selection.Proof.Response), v.powP(m, &selection.Proof.Challenge))
			hash := crypto.Hash1(q, "30")
			hash = crypto.Hash1(q, hash, extendedBaseHash)
			hash = crypto.Hash1(q, hash, k)
			hash = crypto.Hash1(q, hash, encryptedSelection.Ciphertext.Pad)
			hash = crypto.Hash1(q, hash, encryptedSelection.Ciphertext.Data)
			hash = crypto.Hash1(q, hash, a)
			hash = crypto.Hash1(q, hash, b)
			hash = crypto.Hash1(q, hash, m)

			helper.addCheck("(9.C) The challenge is not computed correctly.", selection.Proof.Challenge.Compare(hash))
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
