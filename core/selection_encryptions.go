package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of selection encryptions (Step 4)
	helper := MakeValidationHelper(v.logger, 4, "Correctness of selection encryptions")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.SubmittedBallots
	chunkSize := len(ballots) / 20
	if chunkSize == 0 {
		chunkSize = len(ballots) / 3
	}

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

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	defer helper.wg.Done()

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, ballotSelection := range contest.BallotSelections {
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

				helper.addCheck(step4A1, v.isValidResidue(a))
				helper.addCheck(step4A2, v.isValidResidue(b))
				helper.addCheck(step4A3, v.isValidResidue(a0))
				helper.addCheck(step4A4, v.isValidResidue(b0))
				helper.addCheck(step4A5, v.isValidResidue(a1))
				helper.addCheck(step4A6, v.isValidResidue(b1))
				helper.addCheck(step4B, c.Compare(crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck(step4C1, v.isInRange(c0))
				helper.addCheck(step4C2, v.isInRange(c1))
				helper.addCheck(step4C3, v.isInRange(v0))
				helper.addCheck(step4C4, v.isInRange(v1))
				helper.addCheck(step4D, c.Compare(v.addQ(&c0, &c1)))
				helper.addCheck(step4E, v.powP(v.constants.G, &v0).Compare(v.mulP(&a0, v.powP(&a, &c0))))
				helper.addCheck(step4F, v.powP(v.constants.G, &v1).Compare(v.mulP(&a1, v.powP(&a, &c1))))
				helper.addCheck(step4G, v.powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v0).Compare(v.mulP(&b0, v.powP(&b, &c0))))
				helper.addCheck(step4H, v.mulP(v.powP(v.constants.G, &c1), v.powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v1)).Compare(v.mulP(&b1, v.powP(&b, &c1))))
			}
		}
	}
}
