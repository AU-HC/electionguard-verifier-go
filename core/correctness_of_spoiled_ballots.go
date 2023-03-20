package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateCorrectnessOfSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correctness of spoiled ballots (Step 16)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 16, "Correctness of spoiled ballots")
	xd := make(map[string]map[string]struct{})
	start := time.Now()

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

			} else {
				// TODO: Report *should* be reported
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 16 took: " + time.Since(start).String())
}
