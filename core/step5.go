package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 5, "Selection encryptions are correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.SubmittedBallots
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

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, ballots []schema.SubmittedBallot, er *deserialize.ElectionRecord) {
	defer helper.wg.Done()

	q := &er.ElectionConstants.SmallPrime
	g := &er.ElectionConstants.Generator
	k := &er.CiphertextElectionRecord.ElgamalPublicKey
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				c0 := selection.Proof.ProofZeroChallenge
				v0 := selection.Proof.ProofZeroResponse
				c1 := selection.Proof.ProofOneChallenge
				v1 := selection.Proof.ProofOneResponse
				a := selection.Ciphertext.Pad
				b := selection.Ciphertext.Data

				a0 := v.mulP(v.powP(g, &v0), v.powP(&a, &c0))
				a1 := v.mulP(v.powP(g, &v1), v.powP(&a, &c1))
				b0 := v.mulP(v.powP(k, &v0), v.powP(&b, &c0))
				w1 := v.subQ(&v1, &c1)
				b1 := v.mulP(v.powP(k, w1), v.powP(&b, &c1))
				expectedChallenge := crypto.Hash(q, "21", extendedBaseHash, k, a, b, a0, b0, a1, b1)

				errorString := "(BallotID:" + ballot.Code + ", ContestID:" + contest.ObjectId + ", SelectionID:" + selection.ObjectId + ")"
				helper.addCheck("(5.D) The values alpha and beta are not valid.", v.isValidResidue(a) && v.isValidResidue(b), errorString)
				helper.addCheck("(5.E) The proof values are not valid.", v.isInRange(c0) && v.isInRange(v0) && v.isInRange(c1) && v.isInRange(v1), errorString)
				helper.addCheck("(5.F) The challenge value is not computed correctly.", v.addQ(&c0, &c1).Compare(expectedChallenge), errorString)
			}
		}
	}
}
