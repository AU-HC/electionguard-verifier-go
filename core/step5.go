package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"fmt"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 5, "Selection encryptions are correct")
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
		go v.validateSelectionEncryptionForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, ballots []schema.MockBallot, er *schema.ElectionRecord) {
	defer helper.wg.Done()

	q := er.ElectionConstants.SmallPrime
	g := &er.ElectionConstants.Generator
	k := &er.CiphertextElectionRecord.ElgamalPublicKey
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				alpha := selection.Ciphertext.Pad
				beta := selection.Ciphertext.Data

				toBeHashed := []interface{}{*k, selection.Ciphertext.Pad, selection.Ciphertext.Data}
				computedC := schema.MakeBigIntFromString("0", 10)
				for j, proof := range selection.Proof.Proofs {
					cj := proof.Challenge
					computedC = v.addQ(computedC, &cj)

					vj := proof.Response
					wj := v.subQ(&vj, v.mulQ(schema.MakeBigIntFromString(fmt.Sprintf("%d", j), 10), &cj))

					aj := v.mulP(v.powP(g, &vj), v.powP(&alpha, &cj))
					bj := v.mulP(v.powP(k, wj), v.powP(&beta, &cj))

					toBeHashed = append(toBeHashed, *aj)
					toBeHashed = append(toBeHashed, *bj)

					helper.addCheck("(5.B) The challenge is in the range 0 <= cj < 2^256", v.isInRange(cj))
					helper.addCheck("(5.C) The response vj is in Z_q", v.isInRange(vj))
				}
				c := crypto.HMAC(q, extendedBaseHash, 0x21, toBeHashed...)

				helper.addCheck("(5.A) alpha is in Z_p^r", v.isValidResidue(alpha))
				helper.addCheck("(5.A) beta is in Z_p^r", v.isValidResidue(beta))
				helper.addCheck("(5.D) challenge is computed correctly", computedC.Compare(c))
			}
		}
	}
}
