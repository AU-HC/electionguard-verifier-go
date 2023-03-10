package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of selection encryptions (Step 4)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 4, "Correctness of selection encryptions")

	ballots := er.SubmittedBallots

	// Split the slice of ballots into multiple slices
	chunkSize := len(ballots) / 15
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

				helper.addCheck("(4.A) a is in the set Z_p^r", isValidResidue(a))
				helper.addCheck("(4.A) b is in the set Z_p^r", isValidResidue(b))
				helper.addCheck("(4.A) a0 is in the set Z_p^r", isValidResidue(a0))
				helper.addCheck("(4.A) b0 is in the set Z_p^r", isValidResidue(b0))
				helper.addCheck("(4.A) a1 is in the set Z_p^r", isValidResidue(a1))
				helper.addCheck("(4.A) b1 is in the set Z_p^r", isValidResidue(b1))
				helper.addCheck("(4.B) The challenge value c is computed correctly", c.Compare(crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck("(4.C) c0 is in Zq for", isInRange(c0))
				helper.addCheck("(4.C) c1 is in Zq for", isInRange(c1))
				helper.addCheck("(4.C) v0 is in Zq for", isInRange(v0))
				helper.addCheck("(4.C) v1 is in Zq for", isInRange(v1))
				helper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied", c.Compare(addQ(&c0, &c1)))
				helper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied", powP(v.constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				helper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied", powP(v.constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				helper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied", powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				helper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied", mulP(powP(v.constants.G, &c1), powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))
			}
		}
	}
}
