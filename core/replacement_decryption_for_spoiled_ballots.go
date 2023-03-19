package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateReplacementPartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correct replacement partial decryptions for spoiled ballots (Step 14)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 14, "Correctness of replacement partial decryptions for spoiled ballots")
	start := time.Now()

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				for _, share := range selection.Shares {
					m := share.Share

					if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix {
						product := schema.MakeBigIntFromInt(1)

						for _, part := range share.RecoveredParts {
							coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
							product = v.mulP(product, v.powP(&part.Share, &coefficient))
						}
						if len(share.RecoveredParts) > 0 {
							helper.addCheck("(14.B) Correct missing decryption share", m.Compare(product))
						}
					}
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 14 took: " + time.Since(start).String())
}
