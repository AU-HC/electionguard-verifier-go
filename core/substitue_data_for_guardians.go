package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateSubstituteDataForMissingGuardians(er *deserialize.ElectionRecord) {
	// Validate correctness of substitute data for missing guardians (Step 9)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 9, "Correctness of substitute data for missing guardians")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	// Mapping map to slice
	var contests []schema.ContestTally
	for _, contest := range er.PlaintextTally.Contests {
		contests = append(contests, contest)
	}

	// Split the slice of contests into multiple slices (namely 2)
	chunkSize := len(contests) / 2
	for i := 0; i < len(contests); i += chunkSize {
		end := i + chunkSize

		if end > len(contests) {
			end = len(contests)
		}

		go v.validateSubstituteDataForMissingGuardiansForSlice(helper, contests[i:end], extendedBaseHash)
	}

	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSubstituteDataForMissingGuardiansForSlice(helper *ValidationHelper, contests []schema.ContestTally, extendedBaseHash schema.BigInt) {
	v.wg.Add(1)
	defer v.wg.Done()

	for _, contest := range contests {
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
						m := part.Share

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
}
