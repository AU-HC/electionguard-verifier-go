package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"strconv"
)

func (v *Verifier) validatePartialDecryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of partial decryptions (Step 8)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 8, "Correctness of partial decryptions")

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data

			for k, share := range selection.Shares {
				if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix
					V := share.Proof.Response
					c := share.Proof.Challenge
					ai := share.Proof.Pad
					bi := share.Proof.Data
					m := share.Share

					helper.addCheck("(8.A) The value v is in the set Zq for "+share.ObjectId+" "+strconv.Itoa(k), isInRange(V))
					helper.addCheck("(8.B) The value a is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Pad))
					helper.addCheck("(8.B) The value b is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Data))
					helper.addCheck("(8.C) The challenge is computed correctly "+share.ObjectId+" "+strconv.Itoa(k), c.Compare(crypto.HashElems(extendedBaseHash, A, B, ai, bi, m)))
					helper.addCheck("(8.D) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(v.constants.G, &V).Compare(mulP(&ai, powP(&er.Guardians[k].ElectionPublicKey, &c))))
					helper.addCheck("(8.E) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(&A, &V).Compare(mulP(&bi, powP(&m, &c))))
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
