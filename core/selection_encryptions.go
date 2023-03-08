package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"strconv"
)

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) {
	// Validate correctness of selection encryptions (Step 4)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 4, "Correctness of selection encryptions")

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	elgamalPublicKey := &er.CiphertextElectionRecord.ElgamalPublicKey
	for i, ballot := range er.SubmittedBallots {
		for j, contest := range ballot.Contests {
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
				helper.addCheck("(4.B) The challenge value c is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(crypto.HashElems(extendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck("(4.C) c0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c0))
				helper.addCheck("(4.C) c1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c1))
				helper.addCheck("(4.C) v0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v0))
				helper.addCheck("(4.C) v1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v1))
				helper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(addQ(&c0, &c1)))
				helper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				helper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				helper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(elgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				helper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", mulP(powP(v.constants.G, &c1), powP(elgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
