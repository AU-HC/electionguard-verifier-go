package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"strconv"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of selection encryptions (Step 4)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 4, "Correctness of selection encryptions")

	// extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	// elgamalPublicKey := &er.CiphertextElectionRecord.ElgamalPublicKey
	ballots := er.SubmittedBallots

	// Split the slice of ballots into multiple slices
	chunkSize := len(ballots) / 10
	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		go v.validateSelectionEncryptionForSlice(helper, ballots[i:end], er)
	}

	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	v.wg.Add(1)
	defer v.wg.Done()

	for i, ballot := range ballots {
		for j, contest := range ballot.Contests {
			contestInManifest := getContest(contest.ObjectId, er.Manifest.Contests)
			votesAllowed := contestInManifest.VotesAllowed
			numberOfSelections := 0
			calculatedAHat := schema.MakeBigIntFromInt(1)
			calculatedBHat := schema.MakeBigIntFromInt(1)

			aHat := contest.CiphertextAccumulation.Pad
			bHat := contest.CiphertextAccumulation.Data
			a := contest.Proof.Pad
			b := contest.Proof.Data
			V := contest.Proof.Response

			for k, ballotSelection := range contest.BallotSelections {
				a := ballotSelection.Ciphertext.Pad
				b := ballotSelection.Ciphertext.Data
				a0 := ballotSelection.Proof.ProofZeroPad
				b0 := ballotSelection.Proof.ProofZeroData
				a1 := ballotSelection.Proof.ProofOnePad
				b1 := ballotSelection.Proof.ProofOneData
				c := ballotSelection.Proof.Challenge
				c0 := ballotSelection.Proof.ProofZeroChallenge
				c1 := ballotSelection.Proof.ProofOneChallenge
				v0 := ballotSelection.Proof.ProofZeroResponse
				v1 := ballotSelection.Proof.ProofOneResponse

				helper.addCheck("(4.A) a is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a))
				helper.addCheck("(4.A) b is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b))
				helper.addCheck("(4.A) a0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a0))
				helper.addCheck("(4.A) b0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b0))
				helper.addCheck("(4.A) a1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a1))
				helper.addCheck("(4.A) b1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b1))
				helper.addCheck("(4.B) The challenge value c is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck("(4.C) c0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c0))
				helper.addCheck("(4.C) c1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c1))
				helper.addCheck("(4.C) v0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v0))
				helper.addCheck("(4.C) v1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v1))
				helper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(addQ(&c0, &c1)))
				helper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				helper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				helper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				helper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", mulP(powP(v.constants.G, &c1), powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))

				if ballotSelection.IsPlaceholderSelection {
					numberOfSelections++
				}
				calculatedAHat = mulP(calculatedAHat, &ballotSelection.Ciphertext.Pad)
				calculatedBHat = mulP(calculatedBHat, &ballotSelection.Ciphertext.Data)
			}

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
}
