package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validatePartialDecryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of partial decryptions (Step 8)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 8, "Correctness of partial decryptions")
	start := time.Now()

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data

			for _, share := range selection.Shares {
				if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix
					vi := share.Proof.Response
					c := share.Proof.Challenge
					ai := share.Proof.Pad
					bi := share.Proof.Data
					m := share.Share

					helper.addCheck("(8.A) The value v is in the set Zq for", isInRange(vi))
					helper.addCheck("(8.B) The value a is in the set Zqr for", isValidResidue(share.Proof.Pad))
					helper.addCheck("(8.B) The value b is in the set Zqr for", isValidResidue(share.Proof.Data))
					helper.addCheck("(8.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, A, B, ai, bi, m)))
					helper.addCheck("(8.D) The equation is satisfied", powP(v.constants.G, &vi).Compare(mulP(&ai, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
					helper.addCheck("(8.E) The equation is satisfied", powP(&A, &vi).Compare(mulP(&bi, powP(&m, &c))))
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 8 took: " + time.Since(start).String())
}
