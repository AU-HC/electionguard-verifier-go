package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateCorrectnessOfSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 16, "Correctness of spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	xd := make(map[string]map[string]struct{})
	for _, ballot := range er.SubmittedBallots { // Looping over submitted ballots to have access to ballot state
		if ballot.State != 2 { // State 2 = spoiled
			continue
		}

		spoiledBallot := getSpoiledBallot(ballot.ObjectId, er.SpoiledBallots)
		for _, contest := range spoiledBallot.Contests {
			manifestContest := getContest(contest.ObjectId, er.Manifest.Contests)
			sumOfAllSelections := 0
			xd[contest.ObjectId] = make(map[string]struct{})

			for _, selection := range contest.Selections {
				sumOfAllSelections += selection.Tally
				helper.addCheck(step16A, selection.Tally == 0 || selection.Tally == 1)

				selectionIsNotPlaceholder := !isPlaceholderSelection(ballot, selection.ObjectId)
				if selectionIsNotPlaceholder {
					helper.addCheck(step16D, doesManifestSelectionExist(selection.ObjectId, manifestContest.BallotSelections))
					xd[contest.ObjectId][selection.ObjectId] = struct{}{}
				}
			}
			helper.addCheck(step16B, sumOfAllSelections <= manifestContest.VotesAllowed)
			helper.addCheck(step16C, doesContestExistInManifest(contest.ObjectId, er.Manifest.Contests))
		}
	}

	for _, contest := range er.Manifest.Contests {
		for _, selection := range contest.BallotSelections {
			maps, contestExists := xd[contest.ObjectID]

			if contestExists {
				_, selectionExists := maps[selection.ObjectID]
				helper.addCheck(step16E, selectionExists)

			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
