package core

import "electionguard-verifier-go/deserialize"

func (v *Verifier) validateCorrectnessOfSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validation of correctness of spoiled ballots (Step 16)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 16, "Correctness of spoiled ballots")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			sumOfAllSelections := 0
			for _, selection := range contest.Selections {
				sumOfAllSelections += selection.Tally
				helper.addCheck("(16.A) For each option in the contest, the selection V is either a 0 or a 1", selection.Tally == 0 || selection.Tally == 1)
			}
			helper.addCheck("(16.B) The sum of all selections in the contest is at most the selection limit L for that contest.", sumOfAllSelections <= getContest(contest.ObjectId, er.Manifest.Contests).VotesAllowed)
			// TODO: 16.C -> 16.E
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
