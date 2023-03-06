package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
)

func (v *Verifier) validateSubstituteDataForMissingGuardians(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of substitute data for missing guardians (Step 9)
	helper := MakeValidationHelper(v.logger, "Correctness of substitute data for missing guardians (Step 9)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data
			for _, share := range selection.Shares {
				for _, part := range share.RecoveredParts {
					if part.ObjectId != "" { // TODO: Implement method to check if "Recovered parts" is not nil
						V := part.Proof.Response
						c := part.Proof.Challenge
						a := part.Proof.Pad
						b := part.Proof.Data
						m := part.PartialDecryption

						helper.addCheck("(9.A) The given value v is in Zq", isInRange(V))
						helper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(a))
						helper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(b))
						helper.addCheck("(9.C) The challenge value c is correct", c.Compare(crypto.HashElems(extendedBaseHash, A, B, a, b, m)))
						helper.addCheck("(9.D) The equation is satisfied", powP(v.constants.G, &V).Compare(mulP(&a, powP(&part.RecoveryPublicKey, &c))))
						helper.addCheck("(9.E) The equation is satisfied", powP(&A, &V).Compare(mulP(&b, powP(&m, &c))))
					}
				}
			}
		}
	}

	return helper
}
