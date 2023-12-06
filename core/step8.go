package core

import (
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateBallotAggregation(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 8, "Correctness of ballot aggregation")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// creating set of spoiled ballot codes for easier lookup
	spoiledBallots := map[string]struct{}{}
	for _, spoiledBallot := range er.SpoiledBallots {
		spoiledBallots[spoiledBallot.Name] = struct{}{}
	}

	// verifying step 8
	for _, contest := range er.EncryptedTally.Contests {
		for _, selection := range contest.Selections {
			a := selection.Ciphertext.Pad
			b := selection.Ciphertext.Data

			calculatedA := schema.IntToBigInt(1)
			calculatedB := schema.IntToBigInt(1)

			for _, ballot := range er.SubmittedBallots {
				// check if ballot is spoiled or cast
				_, ballotIsSpoiled := spoiledBallots[ballot.Code]
				if ballotIsSpoiled {
					continue
				}

				// find the correct contest/selection and multiply with the current aggregate
				ballotEncryption := findEncryptionForContestAndSelection(contest.ObjectId, selection.ObjectId, ballot)
				calculatedA = v.mulP(calculatedA, &ballotEncryption.Pad)
				calculatedB = v.mulP(calculatedB, &ballotEncryption.Data)
			}

			errorString := "(ContestID:" + contest.ObjectId + ", SelectionID:" + selection.ObjectId + ")"
			helper.addCheck("(8.A) The ballot aggregation was not correct for A.", a.Compare(calculatedA), errorString)
			helper.addCheck("(8.B) The ballot aggregation was not correct for B.", b.Compare(calculatedB), errorString)
		}
	}

	v.helpers[helper.VerificationStep] = helper
}

func findEncryptionForContestAndSelection(contestID, selectionID string, ballot schema.SubmittedBallot) *schema.Ciphertext {
	for _, contest := range ballot.Contests {
		if contestID == contest.ObjectId {
			for _, selection := range contest.BallotSelections {
				if selectionID == selection.ObjectId {
					return &selection.Ciphertext
				}
			}
		}
	}

	one := *schema.IntToBigInt(1)
	return &schema.Ciphertext{Pad: one, Data: one}
}
