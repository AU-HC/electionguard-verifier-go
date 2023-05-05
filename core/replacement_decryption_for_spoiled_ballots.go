package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateReplacementPartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 14, "Correctness of replacement partial decryptions for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				for _, share := range selection.Shares {
					m := share.Share

					if share.Proof.IsNotEmpty() {
						product := schema.MakeBigIntFromInt(1)

						for _, part := range share.RecoveredParts {
							coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
							product = v.mulP(product, v.powP(&part.Share, &coefficient))
						}
						if len(share.RecoveredParts) > 0 {
							helper.addCheck(step14B, m.Compare(product))
						}
					}
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
