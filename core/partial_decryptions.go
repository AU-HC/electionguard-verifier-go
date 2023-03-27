package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validatePartialDecryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 8, "Correctness of partial decryptions")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

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

					helper.addCheck(step8A, v.isInRange(vi))
					helper.addCheck(step8B1, v.isValidResidue(share.Proof.Pad))
					helper.addCheck(step8B2, v.isValidResidue(share.Proof.Data))
					helper.addCheck(step8C, c.Compare(crypto.HashElems(extendedBaseHash, A, B, ai, bi, m)))
					helper.addCheck(step8D, v.powP(v.constants.G, &vi).Compare(v.mulP(&ai, v.powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
					helper.addCheck(step8E, v.powP(&A, &vi).Compare(v.mulP(&bi, v.powP(&m, &c))))
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
