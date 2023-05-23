package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateVoteLimits(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 5, "Adherence to vote limits")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.SubmittedBallots
	chunkSize := len(ballots) / v.verifierStrategy.getBallotSplitSize()
	if chunkSize == 0 {
		chunkSize = len(ballots) / 3
	}

	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		helper.wg.Add(1)
		go v.validateVoteLimitsForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateVoteLimitsForSlice(helper *ValidationHelper, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	defer helper.wg.Done()

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			contestInManifest := getContest(contest.ObjectId, er.Manifest.Contests)
			votesAllowed := contestInManifest.VotesAllowed
			votesAllowedBigInt := schema.IntToBigInt(votesAllowed)
			numberOfPlaceholderSelections := 0
			calculatedAHat := schema.IntToBigInt(1)
			calculatedBHat := schema.IntToBigInt(1)

			for _, selection := range contest.BallotSelections {
				if selection.IsPlaceholderSelection {
					numberOfPlaceholderSelections++
				}
				calculatedAHat = v.mulP(calculatedAHat, &selection.Ciphertext.Pad)
				calculatedBHat = v.mulP(calculatedBHat, &selection.Ciphertext.Data)
			}

			aHat := contest.CiphertextAccumulation.Pad
			bHat := contest.CiphertextAccumulation.Data
			a := contest.Proof.Pad
			b := contest.Proof.Data
			V := contest.Proof.Response
			c := contest.Proof.Challenge

			computedChallenge := crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, aHat, bHat, a, b)

			helper.addCheck(step5A, votesAllowed == numberOfPlaceholderSelections)
			helper.addCheck(step5B1, aHat.Compare(calculatedAHat))
			helper.addCheck(step5B2, bHat.Compare(calculatedBHat))
			helper.addCheck(step5C, v.isInRange(V))
			helper.addCheck(step5D, v.isValidResidue(a))
			helper.addCheck(step5E, v.isValidResidue(b))
			helper.addCheck(step5F1, c.Compare(computedChallenge))
			helper.addCheck(step5F2, v.powP(v.constants.G, &V).Compare(v.mulP(&a, v.powP(&aHat, &c))))
			helper.addCheck(step5G, v.mulP(v.powP(v.constants.G, v.mulP(votesAllowedBigInt, &c)), v.powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &V)).Compare(v.mulP(&b, v.powP(&bHat, &c))))
		}
	}
}
