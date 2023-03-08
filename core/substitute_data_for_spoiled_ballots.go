package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateSubstituteDataForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validating correctness of substitute data for spoiled ballots (Step 13)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 13, "Correctness of substitute data for spoiled ballots")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) {
						m := share.Share
						product := schema.MakeBigIntFromInt(1)

						for _, part := range share.RecoveredParts {
							coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
							product = mulP(product, powP(&part.Share, &coefficient))
						}
						if len(share.RecoveredParts) > 0 {
							helper.addCheck("(14.B) Correct missing decryption share", m.Compare(product))
						}
					}
				}
			}
		}
	}

	v.helpers[helper.verificationStep] = helper
}
