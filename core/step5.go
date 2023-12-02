package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 5, "Selection encryptions are correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	q := &er.ElectionConstants.SmallPrime
	g := &er.ElectionConstants.Generator
	k := &er.CiphertextElectionRecord.ElgamalPublicKey
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SubmittedBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				c0 := selection.Proof.ProofZeroChallenge
				v0 := selection.Proof.ProofZeroResponse
				c1 := selection.Proof.ProofOneChallenge
				v1 := selection.Proof.ProofOneResponse
				a := selection.Ciphertext.Pad
				b := selection.Ciphertext.Data

				a0 := v.mulP(v.powP(g, &v0), v.powP(&a, &c0))
				a1 := v.mulP(v.powP(g, &v1), v.powP(&a, &c1))
				b0 := v.mulP(v.powP(k, &v0), v.powP(&b, &c0))
				w1 := v.subQ(&v1, &c1)
				b1 := v.mulP(v.powP(k, w1), v.powP(&b, &c1))
				expectedChallenge := crypto.Hash1(q, "21", extendedBaseHash, k, a, b, a0, b0, a1, b1)

				helper.addCheck("(5.D) TODO.", v.isValidResidue(a) && v.isValidResidue(b))
				helper.addCheck("(5.E) TODO.", v.isInRange(c0) && v.isInRange(v0) && v.isInRange(c1) && v.isInRange(v1))
				helper.addCheck("(5.F) TODO.", v.addQ(&c0, &c1).Compare(expectedChallenge))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
