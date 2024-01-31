package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"fmt"
	"time"
)

func (v *Verifier) validateAdherenceToVoteLimits(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 6, "Adherence to vote limits")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.MockBallots
	chunkSize := v.verifierStrategy.getBallotChunkSize(len(ballots))
	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		helper.wg.Add(1)
		go v.validateAdherenceToVoteLimitsForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateAdherenceToVoteLimitsForSlice(helper *ValidationHelper, ballots []schema.MockBallot, er *schema.ElectionRecord) {
	defer helper.wg.Done()

	q := er.ElectionConstants.SmallPrime
	g := &er.ElectionConstants.Generator
	k := &er.CiphertextElectionRecord.ElgamalPublicKey
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			alphaHat, betaHat := v.computeContestTotal(contest.BallotSelections)

			toBeHashed := []interface{}{k, alphaHat, betaHat}
			computedC := schema.MakeBigIntFromString("0", 10)
			for j, proof := range contest.Proof.Proofs {
				cj := proof.Challenge
				computedC = v.addQ(computedC, &cj)

				vj := proof.Response
				wj := v.subQ(&vj, v.mulQ(schema.MakeBigIntFromString(fmt.Sprintf("%d", j), 10), &cj))

				aj := v.mulP(v.powP(g, &vj), v.powP(alphaHat, &cj))
				bj := v.mulP(v.powP(k, wj), v.powP(betaHat, &cj))

				toBeHashed = append(toBeHashed, *aj)
				toBeHashed = append(toBeHashed, *bj)

				// TODO: Check A is already done in step5? also check 6.B

				helper.addCheck("(6.B) The challenge is in the range 0 <= cj < 2^256", v.isInRange(cj))
				helper.addCheck("(6.C) The response vj is in Z_q", v.isInRange(vj))
			}
			c := crypto.HMAC(q, extendedBaseHash, 0x21, toBeHashed...)
			helper.addCheck("(6.D) Challenge is computed correctly", c.Compare(computedC))
		}
	}
}

func (v *Verifier) computeContestTotal(ballotSelections []schema.MockBallotSelection) (*schema.BigInt, *schema.BigInt) {
	alphaHat := schema.MakeBigIntFromString("1", 10)
	betaHat := schema.MakeBigIntFromString("1", 10)

	for _, selection := range ballotSelections {
		alphaHat = v.mulP(alphaHat, &selection.Ciphertext.Pad)
		betaHat = v.mulP(betaHat, &selection.Ciphertext.Data)
	}

	return alphaHat, betaHat
}
