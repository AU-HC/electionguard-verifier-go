package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"strconv"
)

func (v *Verifier) validateVoteLimits(er *deserialize.ElectionRecord) {
	// Validate adherence to vote limits (Step 5)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 5, "Adherence to vote limits")

	for i, ballot := range er.SubmittedBallots {
		for j, contest := range ballot.Contests {
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

			helper.addCheck("(5.A) The number of placeholder positions matches the selection limit ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", votesAllowed == numberOfSelections)
			helper.addCheck("(5.B) The a hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", aHat.Compare(calculatedAHat))
			helper.addCheck("(5.B) The b hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", bHat.Compare(calculatedBHat))
			helper.addCheck("(5.C) The given value V is in Zq ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isInRange(V))
			helper.addCheck("(5.D) The given value a are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Pad))
			helper.addCheck("(5.D) The given values b are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Data))
			helper.addCheck("(5.E) The challenge value is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", contest.Proof.Challenge.Compare(c))
			helper.addCheck("(5.F) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationFLeft.Compare(equationFRight))
			helper.addCheck("(5.G) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationGLeft.Compare(equationGRight))
		}
	}

	v.helpers[helper.verificationStep] = helper
}
