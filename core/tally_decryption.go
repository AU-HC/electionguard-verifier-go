package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateTallyDecryption(er *deserialize.ElectionRecord) {
	// Validate correctness of tally decryption (Step 11)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 11, "Correct decryption of tallies")
	contestMap := make(map[string]map[string]struct{})
	ballotContestMap := make(map[string]struct{})
	start := time.Now()

	for _, contest := range er.PlaintextTally.Contests {
		contestMap[contest.ObjectId] = make(map[string]struct{})
		helper.addCheck(step11C, contains(er.Manifest.Contests, contest.ObjectId))

		for _, selection := range contest.Selections {
			contestMap[contest.ObjectId][selection.ObjectId] = struct{}{}

			b := selection.Message.Data
			m := selection.Value
			t := schema.MakeBigIntFromInt(selection.Tally)

			mi := schema.MakeBigIntFromString("1", 10)
			for _, share := range selection.Shares {
				mi = v.mulP(mi, &share.Share)
			}

			helper.addCheck(step11A, b.Compare(v.mulP(&m, mi)))
			helper.addCheck(step11B, m.Compare(v.powP(v.constants.G, t)))
			helper.addCheck(step11D, doesContestContainSelection(er.Manifest.Contests, contest.ObjectId, selection.ObjectId))
		}
	}

	for _, ballot := range er.SubmittedBallots {
		for _, contest := range ballot.Contests {
			ballotContestMap[contest.ObjectId] = struct{}{}
		}
	}

	// TODO: This check does not take placeholder options into account
	// as that information is contained in submitted ballots
	for _, contest := range er.Manifest.Contests {
		contestSelections, ok := contestMap[contest.ObjectID]
		if ok {
			for _, selection := range contest.BallotSelections {
				_, ok = contestSelections[selection.ObjectID]
				helper.addCheck(step11E, ok)
			}
		} else {
			// error *should* already be logged
		}

		_, ok = ballotContestMap[contest.ObjectID]
		helper.addCheck(step11F, ok)
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 11 took: " + time.Since(start).String())
}
