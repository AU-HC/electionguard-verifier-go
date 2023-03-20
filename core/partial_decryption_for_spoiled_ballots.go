package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validatePartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 12, "Correctness of partial decryption for spoiled ballots")
	start := time.Now()

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) { // Comparing with zero, will need better way of determining this TODO: Fix
						m := share.Share
						a := share.Proof.Pad
						b := share.Proof.Data
						c := share.Proof.Challenge
						V := share.Proof.Response

						helper.addCheck("(12.A) The given value v is in the set Zq", v.isInRange(V))
						helper.addCheck("(12.B) The given value a is in the set Zpr", v.isValidResidue(a))
						helper.addCheck("(12.B) The given value b is in the set Zpr", v.isValidResidue(b))
						helper.addCheck("(12.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, a, b, m)))
						helper.addCheck("(12.D) The equation is satisfied", v.powP(v.constants.G, &V).Compare(v.mulP(&a, v.powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
						helper.addCheck("(12.E) The equation is satisfied", v.powP(&alpha, &V).Compare(v.mulP(&b, v.powP(&m, &c))))
					}
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 12 took: " + time.Since(start).String())
}
