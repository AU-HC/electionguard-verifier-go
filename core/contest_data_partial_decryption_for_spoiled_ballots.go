package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
)

func (v *Verifier) validateContestDataPartialDecryptionsForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 17, "Correctness of contest data partial decryptions for spoiled ballots")
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

				helper.addCheck("(17.A) The given value v is in the set Zq", isInRange(vi))
				helper.addCheck("(17.B) The given value a is in the set Zqr", isValidResidue(ai))
				helper.addCheck("(17.B) The given value b is in the set Zqr", isValidResidue(bi))
				helper.addCheck("(17.C) The challenge is correctly computed", ci.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, ai, bi, mi)))
				helper.addCheck("(17.D) The equation is satisfied", powP(v.constants.G, &vi).Compare(mulP(&ai, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &ci))))
				helper.addCheck("(17.E) The equation is satisfied", powP(&c0, &vi).Compare(mulP(&bi, powP(&mi, &ci))))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
