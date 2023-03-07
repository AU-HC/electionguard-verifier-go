package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateContestReplacementDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correctness of contest replacement decryptions for spoiled ballots (Step 19)
	helper := MakeValidationHelper(v.logger, "Correctness of contest replacement decryptions for spoiled ballots (Step 19)")
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				product := schema.MakeBigIntFromInt(1)
				for _, part := range share.RecoveredParts {
					m := part.Share

					coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
					product = mulP(product, powP(&m, &coefficient))
				}
				helper.addCheck("(19.A) The equation is satisfied", mi.Compare(product))
			}
		}
	}

	return helper
}
