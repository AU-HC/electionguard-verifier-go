package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateDecryptionOfSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correct decryption of spoiled ballots (Step 15)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 15, "Correct decryption of spoiled ballots")
	start := time.Now()

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			sumOfAllSelections := 0
			for _, selection := range contest.Selections {
				beta := selection.Message.Data
				m := selection.Value
				V := schema.MakeBigIntFromInt(selection.Tally)
				mi := schema.MakeBigIntFromInt(1)
				sumOfAllSelections += selection.Tally
				for _, share := range selection.Shares {
					mi = v.mulP(mi, &share.Share)
				}

				helper.addCheck(step15A, beta.Compare(v.mulP(&m, mi)))
				helper.addCheck(step15B, m.Compare(v.powP(v.constants.G, V)))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 15 took: " + time.Since(start).String())
}
