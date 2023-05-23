package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateContestReplacementDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 19, "Correctness of contest replacement decryptions for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				product := schema.IntToBigInt(1)
				for _, part := range share.RecoveredParts {
					m := part.Share

					coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
					product = v.mulP(product, v.powP(&m, &coefficient))
				}
				helper.addCheck(step19A, mi.Compare(product))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
