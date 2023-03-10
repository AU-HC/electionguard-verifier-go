package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateBallotAggregation(er *deserialize.ElectionRecord) {
	// Validate correctness of ballot aggregation (Step 7)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 7, "Correctness of ballot aggregation")

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			a := schema.MakeBigIntFromInt(1)
			b := schema.MakeBigIntFromInt(1)
			for _, ballot := range er.SubmittedBallots {
				ballotWasCast := ballot.State == 1
				if ballotWasCast {
					ciphertextSelection := getSelection(ballot, contest.ObjectId, selection.ObjectId)
					a = mulP(a, &ciphertextSelection.Pad)
					b = mulP(b, &ciphertextSelection.Data)
				}
			}
			A := selection.Message.Pad
			B := selection.Message.Data
			helper.addCheck("(7.A) A is calculated correctly", A.Compare(a))
			helper.addCheck("(7.B) B is calculated correctly", B.Compare(b))
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
