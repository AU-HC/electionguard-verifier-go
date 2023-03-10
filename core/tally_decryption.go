package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateTallyDecryption(er *deserialize.ElectionRecord) {
	// Validate correctness of tally decryption (Step 11)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 11, "Correct decryption of tallies")

	for _, contest := range er.PlaintextTally.Contests {
		helper.addCheck("Tally label exists in election manifest", contains(er.Manifest.Contests, contest.ObjectId))
		// TODO: Check 11.C to 11.F
		for _, selection := range contest.Selections {
			b := selection.Message.Data
			mi := schema.MakeBigIntFromString("1", 10)
			m := selection.Value
			t := schema.MakeBigIntFromInt(selection.Tally)
			for _, share := range selection.Shares {
				mi = mulP(mi, &share.Share)
			}
			helper.addCheck("(11.A) The equation is satisfied", b.Compare(mulP(&m, mi)))
			helper.addCheck("(11.B) The equation is satisfied", m.Compare(powP(v.constants.G, t)))
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
