package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateBallotAggregation(er *deserialize.ElectionRecord) {
	// Validate correctness of ballot aggregation (Step 7)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 7, "Correctness of ballot aggregation")
	start := time.Now()

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			a := schema.MakeBigIntFromInt(1)
			b := schema.MakeBigIntFromInt(1)
			for _, ballot := range er.SubmittedBallots {
				ballotWasCast := ballot.State == 1
				if ballotWasCast {
					ciphertextSelection := getSelection(ballot, contest.ObjectId, selection.ObjectId)
					a = v.mulP(a, &ciphertextSelection.Pad)
					b = v.mulP(b, &ciphertextSelection.Data)
				}
			}
			A := selection.Message.Pad
			B := selection.Message.Data
			helper.addCheck(step7A, A.Compare(a))
			helper.addCheck(step7B, B.Compare(b))
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 7 took: " + time.Since(start).String())
}
