package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateCorrectnessOfSpoiledBallotDecryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 13, "Validation of correct decryption of spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	k := er.CiphertextElectionRecord.ElgamalPublicKey
	contests := er.Manifest.Contests

	for _, spoiledBallot := range er.SpoiledBallots {
		tallyContests := spoiledBallot.Contests

		// checking if all selections in the manifest exists in the plaintext tally
		for _, contest := range contests {
			contestID := contest.ObjectID
			tallyContest, contestExistsInTally := tallyContests[contestID]

			if contestExistsInTally {
				for _, selection := range contest.BallotSelections {
					selectionID := selection.ObjectID
					errorString := "(SpoiledBallotID:" + spoiledBallot.Name + ", ContestID:" + contestID + ", SelectionID:" + selectionID + ")"

					_, selectionExistsForContest := tallyContest.Selections[selectionID]
					helper.addCheck("(13.F) The option text label does not occur in a spoiled ballot.", selectionExistsForContest, errorString)
				}
			}
		}

		// checking if every contest selection in the spoiled ballot is in the manifest, and validating the "tally" for the spoiled ballot
		for contestID, tallyContest := range tallyContests {
			errorString := "(SpoiledBallotID:" + spoiledBallot.Name + ", ContestID:" + contestID + ")"

			contest, contestExistsInManifest := isContestIDInManifest(contestID, contests)
			helper.addCheck("(13.D) ContestID missing in manifest contests.", contestExistsInManifest, errorString)
			if !contestExistsInManifest {
				continue // if the key is not present simply record the error and continue to next contest
			}

			selectionSum := 0
			selections := contest.BallotSelections
			for selectionID, tallySelection := range tallyContest.Selections {
				errorString = "(SpoiledBallotID:" + spoiledBallot.Name + ", ContestID:" + contestID + ", SelectionID:" + selectionID + ")"

				selectionExistsInManifest := isSelectionIDInManifest(selectionID, selections)
				helper.addCheck("(13.E) SelectionID missing in manifest selections.", selectionExistsInManifest, errorString)
				if !selectionExistsInManifest {
					continue // if the key is not present simply record the error and continue to next contest
				}

				selectionSum += tallySelection.Tally
				S := tallySelection.Value
				sigma := schema.IntToBigInt(tallySelection.Tally)

				helper.addCheck("(13.A) The tally value is incorrectly computed.", S.Compare(v.powP(&k, sigma)), errorString)
				helper.addCheck("(13.B) The selection is not a valid value.", tallySelection.Tally == 0 || tallySelection.Tally == 1, errorString)
			}

			errorString = "(SpoiledBallotID:" + spoiledBallot.Name + ", ContestID:" + contestID + ")"
			helper.addCheck("(13.C) The sum of all selections is larger than the selection limit.", selectionSum <= contest.VotesAllowed, errorString)
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
