package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateVoteLimits(er *deserialize.ElectionRecord) {
	// Validate adherence to vote limits (Step 5)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 5, "Adherence to vote limits")

	// Split the slice of ballots into multiple slices
	ballots := er.SubmittedBallots
	chunkSize := len(ballots) / 15
	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		go v.validateVoteLimitsForSlice(helper, ballots[i:end], er)
	}

	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateVoteLimitsForSlice(helper *ValidationHelper, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	v.wg.Add(1)
	defer v.wg.Done()

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			contestInManifest := getContest(contest.ObjectId, er.Manifest.Contests)
			votesAllowed := contestInManifest.VotesAllowed
			numberOfSelections := 0
			calculatedAHat := schema.MakeBigIntFromInt(1)
			calculatedBHat := schema.MakeBigIntFromInt(1)

			for _, selection := range contest.BallotSelections {
				if selection.IsPlaceholderSelection {
					numberOfSelections++
				}
				calculatedAHat = mulP(calculatedAHat, &selection.Ciphertext.Pad)
				calculatedBHat = mulP(calculatedBHat, &selection.Ciphertext.Data)
			}

			aHat := contest.CiphertextAccumulation.Pad
			bHat := contest.CiphertextAccumulation.Data
			a := contest.Proof.Pad
			b := contest.Proof.Data
			V := contest.Proof.Response

			c := crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, aHat, bHat, a, b)
			equationFLeft := powP(v.constants.G, &V)
			equationFRight := mulP(&a, powP(&aHat, c))
			equationGLeft := mulP(powP(v.constants.G, mulP(schema.MakeBigIntFromInt(votesAllowed), c)), powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &V))
			equationGRight := mulP(&b, powP(&bHat, c))

			helper.addCheck("(5.A) The number of placeholder positions matches the selection limit", votesAllowed == numberOfSelections)
			helper.addCheck("(5.B) The a hat is computed correctly", aHat.Compare(calculatedAHat))
			helper.addCheck("(5.B) The b hat is computed correctly", bHat.Compare(calculatedBHat))
			helper.addCheck("(5.C) The given value V is in Zq", isInRange(V))
			helper.addCheck("(5.D) The given value a are in Zp^r", isValidResidue(contest.Proof.Pad))
			helper.addCheck("(5.D) The given values b are in Zp^r", isValidResidue(contest.Proof.Data))
			helper.addCheck("(5.E) The challenge value is correctly computed", contest.Proof.Challenge.Compare(c))
			helper.addCheck("(5.F) The equation is satisfied", equationFLeft.Compare(equationFRight))
			helper.addCheck("(5.G) The equation is satisfied", equationGLeft.Compare(equationGRight))
		}
	}

}
