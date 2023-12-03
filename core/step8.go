package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateBallotAggregation(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 8, "Correctness of ballot aggregation")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// TODO: check only contests/selections in manifest file
	for _, contest := range er.EncryptedTally.Contests {
		for _, selection := range contest.Selections {
			a := selection.Ciphertext.Pad
			b := selection.Ciphertext.Data

			calculatedA := schema.MakeBigIntFromString("1", 10)
			calculatedB := schema.MakeBigIntFromString("1", 10)

			for _, ballot := range er.SubmittedBallots {
				// check if ballot is spoiled or cast
				ballotIsSpoiled := isBallotSpoiled(ballot, er)
				if ballotIsSpoiled {
					continue
				}

				// find the correct contest/selection and multiply with the current aggregate
				ballotPad, ballotData := findEncryptionForContestAndSelection(contest.ObjectId, selection.ObjectId, ballot)
				calculatedA = v.mulP(calculatedA, ballotPad)
				calculatedB = v.mulP(calculatedB, ballotData)
			}

			helper.addCheck("(8.A) The ballot aggregation was not correct for A.", a.Compare(calculatedA))
			helper.addCheck("(8.B) The ballot aggregation was not correct for B.", b.Compare(calculatedB))
		}
	}

	v.helpers[helper.VerificationStep] = helper
}

// TODO: handle error if contest/selection ids not found
func findEncryptionForContestAndSelection(contestID, selectionID string, ballot schema.SubmittedBallot) (*schema.BigInt, *schema.BigInt) {
	for _, contest := range ballot.Contests {
		if contestID == contest.ObjectId {
			for _, selection := range contest.BallotSelections {
				if selectionID == selection.ObjectId {
					return &selection.Ciphertext.Pad, &selection.Ciphertext.Data
				}
			}
		}
	}

	one := schema.MakeBigIntFromString("1", 10)
	return one, one
}

// TODO: create set of spoiledBallot names for quicker lookup
func isBallotSpoiled(ballot schema.SubmittedBallot, er *deserialize.ElectionRecord) bool {
	for _, spoiledBallot := range er.SpoiledBallots {
		if ballot.Code == spoiledBallot.Name {
			return true
		}
	}
	return false
}
