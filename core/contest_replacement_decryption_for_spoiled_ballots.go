package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateContestReplacementDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correctness of contest replacement decryptions for spoiled ballots (Step 19)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 19, "Correctness of contest replacement decryptions for spoiled ballots")

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

	v.helpers[helper.VerificationStep] = helper
}
