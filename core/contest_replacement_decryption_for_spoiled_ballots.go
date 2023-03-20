package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateContestReplacementDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correctness of contest replacement decryptions for spoiled ballots (Step 19)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 19, "Correctness of contest replacement decryptions for spoiled ballots")
	start := time.Now()

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				product := schema.MakeBigIntFromInt(1)
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
	v.logger.Info("Validation of step 19 took: " + time.Since(start).String())
}
