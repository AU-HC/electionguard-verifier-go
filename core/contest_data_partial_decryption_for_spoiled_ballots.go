package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateContestDataPartialDecryptionsForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 17, "Correctness of contest data partial decryptions for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			c0 := contest.ContestData.Ciphertext.Generator
			c1 := contest.ContestData.Ciphertext.EncryptedMessage
			c2 := contest.ContestData.Ciphertext.MessageAuthenticationCode

			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				ai := share.Proof.Pad
				bi := share.Proof.Data
				ci := share.Proof.Challenge
				vi := share.Proof.Response

				helper.addCheck(step17A, v.isInRange(vi))
				helper.addCheck(step17B1, v.isValidResidue(ai))
				helper.addCheck(step17B2, v.isValidResidue(bi))
				helper.addCheck(step17C, ci.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, ai, bi, mi)))
				helper.addCheck(step17D, v.powP(v.constants.G, &vi).Compare(v.mulP(&ai, v.powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &ci))))
				helper.addCheck(step17E, v.powP(&c0, &vi).Compare(v.mulP(&bi, v.powP(&mi, &ci))))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
