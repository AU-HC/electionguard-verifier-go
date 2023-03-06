package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateReplacementPartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correct replacement partial decryptions for spoiled ballots (Step 14)
	helper := MakeValidationHelper(v.logger, "Correctness of replacement partial decryptions for spoiled ballots (Step 14)")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				for _, share := range selection.Shares {
					m := share.Share

					if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix {
						product := schema.MakeBigIntFromInt(1)

						for _, part := range share.RecoveredParts {
							coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
							product = mulP(product, powP(&part.PartialDecryption, &coefficient))
						}
						if len(share.RecoveredParts) > 0 {
							helper.addCheck("(14.B) Correct missing decryption share", m.Compare(product))
						}
					}
				}
			}
		}
	}

	return helper
}
