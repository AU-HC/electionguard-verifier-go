package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateCorrectnessOfTallyDecryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 10, "Correct decryption of tallies")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	k := er.CiphertextElectionRecord.ElgamalPublicKey
	contests := er.Manifest.Contests
	tallyContests := er.PlaintextTally.Contests

	// checking if all selections in the manifest exists in the plaintext tally
	for _, contest := range contests {
		contestID := contest.ObjectID
		tallyContest, contestExistsInTally := tallyContests[contestID]

		if contestExistsInTally {
			for _, selection := range contest.BallotSelections {
				selectionID := selection.ObjectID
				errorString := "(ContestID:" + contestID + ", SelectionID:" + selectionID + ")"

				_, selectionExistsForContest := tallyContest.Selections[selectionID]
				helper.addCheck("(10.D) Tally missing selectionID in contestID.", selectionExistsForContest, errorString)
			}
		}
	}

	// checking if all contests in all ballots exists in the tally
	submittedBallotContests := allContests(er)
	for ballotContest, _ := range submittedBallotContests {
		errorString := "(ContestID:" + ballotContest + ")"
		_, contestIDExistsInTally := tallyContests[ballotContest]
		helper.addCheck("(10.E) Ballot contest has missing contest in tally.", contestIDExistsInTally, errorString)
	}

	// checking if every contest selection in the tally is in the manifest, and validating the tally
	for _, tallyContest := range tallyContests {
		contestID := tallyContest.ObjectId
		errorString := "(ContestID:" + contestID + ")"

		contest, contestExistsInManifest := isContestIDInManifest(contestID, contests)
		helper.addCheck("(10.B) ContestID missing in manifest contests.", contestExistsInManifest, errorString)
		if !contestExistsInManifest {
			continue // if the key is not present simply record the error and continue to next contest
		}

		selections := contest.BallotSelections
		for _, tallySelection := range tallyContest.Selections {
			errorString = "(ContestID:" + contestID + ", SelectionID:" + tallySelection.ObjectId + ")"

			selectionExistsInManifest := isSelectionIDInManifest(tallySelection.ObjectId, selections)
			helper.addCheck("(10.C) SelectionID missing in manifest selections.", selectionExistsInManifest, errorString)
			if !selectionExistsInManifest {
				continue // if the key is not present simply record the error and continue to next contest
			}

			T := tallySelection.Value
			t := schema.IntToBigInt(tallySelection.Tally)

			helper.addCheck("(10.A) The tally value is incorrectly computed", T.Compare(v.powP(&k, t)), errorString)
		}
	}

	v.helpers[helper.VerificationStep] = helper
}

func isSelectionIDInManifest(selectionID string, selections []schema.ManifestBallotSelection) bool {
	for _, selection := range selections {
		if selectionID == selection.ObjectID {
			return true
		}
	}
	return false
}

func isContestIDInManifest(contestID string, contests []schema.ManifestContest) (schema.ManifestContest, bool) {
	for _, contest := range contests {
		if contestID == contest.ObjectID {
			return contest, true
		}
	}
	return schema.ManifestContest{}, false
}

func allContests(er *deserialize.ElectionRecord) map[string]struct{} {
	result := make(map[string]struct{})

	for _, ballot := range er.SubmittedBallots {
		for _, contest := range ballot.Contests {
			result[contest.ObjectId] = struct{}{}
		}
	}

	return result
}
