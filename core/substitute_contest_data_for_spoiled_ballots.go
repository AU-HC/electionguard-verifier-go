package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateSubstituteContestDataForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 18, "Correctness of substitute contest data for spoiled ballots")
	start := time.Now()

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

					helper.addCheck(step18A, v.isInRange(V))
					helper.addCheck(step18B1, v.isValidResidue(a))
					helper.addCheck(step18B2, v.isValidResidue(b))
					helper.addCheck(step18C, c.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, a, b, m)))
					helper.addCheck(step18D, v.powP(v.constants.G, &V).Compare(v.mulP(&a, v.powP(&er.Guardians[k].ElectionPublicKey, &c))))
					helper.addCheck(step18E, v.powP(&c0, &V).Compare(v.mulP(&b, v.powP(&m, &c))))
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 18 took: " + time.Since(start).String())
}
