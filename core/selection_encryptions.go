package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"sync"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of selection encryptions (Step 4)
	defer v.wg.Done()
	start := time.Now()
	helper := MakeValidationHelper(v.logger, 4, "Correctness of selection encryptions")
	var step4Wg sync.WaitGroup

	ballots := er.SubmittedBallots

	// Split the slice of ballots into multiple slices
	chunkSize := len(ballots) / 20
	if chunkSize == 0 {
		chunkSize = len(ballots) / 3
	}

	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		step4Wg.Add(1)
		go v.validateSelectionEncryptionForSlice(helper, &step4Wg, ballots[i:end], er)
	}

	step4Wg.Wait()
	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 4 took: " + time.Since(start).String())
}

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, wg *sync.WaitGroup, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	defer wg.Done()

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

				helper.addCheck("(4.A) a is in the set Z_p^r", v.isValidResidue(a))
				helper.addCheck("(4.A) b is in the set Z_p^r", v.isValidResidue(b))
				helper.addCheck("(4.A) a0 is in the set Z_p^r", v.isValidResidue(a0))
				helper.addCheck("(4.A) b0 is in the set Z_p^r", v.isValidResidue(b0))
				helper.addCheck("(4.A) a1 is in the set Z_p^r", v.isValidResidue(a1))
				helper.addCheck("(4.A) b1 is in the set Z_p^r", v.isValidResidue(b1))
				helper.addCheck("(4.B) The challenge value c is computed correctly", c.Compare(crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck("(4.C) c0 is in Zq for", v.isInRange(c0))
				helper.addCheck("(4.C) c1 is in Zq for", v.isInRange(c1))
				helper.addCheck("(4.C) v0 is in Zq for", v.isInRange(v0))
				helper.addCheck("(4.C) v1 is in Zq for", v.isInRange(v1))
				helper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied", c.Compare(v.addQ(&c0, &c1)))
				helper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied", v.powP(v.constants.G, &v0).Compare(v.mulP(&a0, v.powP(&a, &c0))))
				helper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied", v.powP(v.constants.G, &v1).Compare(v.mulP(&a1, v.powP(&a, &c1))))
				helper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied", v.powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v0).Compare(v.mulP(&b0, v.powP(&b, &c0))))
				helper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied", v.mulP(v.powP(v.constants.G, &c1), v.powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v1)).Compare(v.mulP(&b1, v.powP(&b, &c1))))
			}
		}
	}
}
