package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateTallyDecryption(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 11, "Correct decryption of tallies")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	contestMap := make(map[string]map[string]struct{})
	ballotContestMap := make(map[string]struct{})

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

	v.helpers[helper.VerificationStep] = helper
}
