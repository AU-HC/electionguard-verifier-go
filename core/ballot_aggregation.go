package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateBallotAggregation(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 7, "Correctness of ballot aggregation")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			a := schema.IntToBigInt(1)
			b := schema.IntToBigInt(1)
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
}
