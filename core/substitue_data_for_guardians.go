package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateSubstituteDataForMissingGuardians(er *deserialize.ElectionRecord) {
	// Validate correctness of substitute data for missing guardians (Step 9)
	helper := MakeValidationHelper(v.logger, 9, "Correctness of substitute data for missing guardians")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	// Mapping map to slice
	var contests []schema.ContestTally
	for _, contest := range er.PlaintextTally.Contests {
		contests = append(contests, contest)
	}

	// Split the slice of contests into multiple slices (namely 1 or 2)
	chunkSize := 1
	if len(contests) > v.verifierStrategy.getContestSplitSize() {
		chunkSize = len(contests) / v.verifierStrategy.getContestSplitSize()
	}

	for i := 0; i < len(contests); i += chunkSize {
		end := i + chunkSize

		if end > len(contests) {
			end = len(contests)
		}

		helper.wg.Add(1)
		go v.validateSubstituteDataForMissingGuardiansForSlice(helper, contests[i:end], extendedBaseHash)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSubstituteDataForMissingGuardiansForSlice(helper *ValidationHelper, contests []schema.ContestTally, extendedBaseHash schema.BigInt) {
	defer helper.wg.Done()

	for _, contest := range contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data
			for _, share := range selection.Shares {
				for _, part := range share.RecoveredParts {
					if part.IsNotEmpty() {
						V := part.Proof.Response
						c := part.Proof.Challenge
						a := part.Proof.Pad
						b := part.Proof.Data
						m := part.Share

						helper.addCheck(step9A, v.isInRange(V))
						helper.addCheck(step9B1, v.isValidResidue(a))
						helper.addCheck(step9B2, v.isValidResidue(b))
						helper.addCheck(step9C, c.Compare(crypto.HashElems(extendedBaseHash, A, B, a, b, m)))
						helper.addCheck(step9D, v.powP(v.constants.G, &V).Compare(v.mulP(&a, v.powP(&part.RecoveryPublicKey, &c))))
						helper.addCheck(step9E, v.powP(&A, &V).Compare(v.mulP(&b, v.powP(&m, &c))))
					}
				}
			}
		}
	}
}
