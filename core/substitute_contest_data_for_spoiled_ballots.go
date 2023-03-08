package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
)

func (v *Verifier) validateSubstituteContestDataForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 18, "Correctness of substitute contest data for spoiled ballots")

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			c0 := contest.ContestData.Ciphertext.Generator
			c1 := contest.ContestData.Ciphertext.EncryptedMessage
			c2 := contest.ContestData.Ciphertext.MessageAuthenticationCode

			for k, share := range contest.ContestData.Shares {
				for _, part := range share.RecoveredParts {
					V := part.Proof.Response
					c := part.Proof.Challenge
					a := part.Proof.Pad
					b := part.Proof.Data
					m := part.Share

					helper.addCheck("(18.A) The value v is in the set Zq", isInRange(V))
					helper.addCheck("(18.B) The value a is in the set Zqr", isValidResidue(a))
					helper.addCheck("(18.B) The value b is in the set Zqr", isValidResidue(b))
					helper.addCheck("(18.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, a, b, m)))
					helper.addCheck("(18.D) The equation is satisfied", powP(v.constants.G, &V).Compare(mulP(&a, powP(&er.Guardians[k].ElectionPublicKey, &c))))
					helper.addCheck("(18.E) The equation is satisfied", powP(&c0, &V).Compare(mulP(&b, powP(&m, &c))))
				}
			}
		}
	}

	v.helpers[helper.verificationStep] = helper
}
