package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validatePartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 12, "Correctness of partial decryption for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	ballots := er.SpoiledBallots
	chunkSize := 1
	if len(ballots) > v.verifierStrategy.getBallotSplitSize() {
		chunkSize = len(ballots) / v.verifierStrategy.getBallotSplitSize()
	}

	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		helper.wg.Add(1)
		go v.validatePartialDecryptionForSpoiledBallotsForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validatePartialDecryptionForSpoiledBallotsForSlice(helper *ValidationHelper, spoiledBallots []schema.SpoiledBallot, er *deserialize.ElectionRecord) {
	defer helper.wg.Done()

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range spoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data
				for _, share := range selection.Shares {
					if share.Proof.IsNotEmpty() {
						m := share.Share
						a := share.Proof.Pad
						b := share.Proof.Data
						c := share.Proof.Challenge
						V := share.Proof.Response

						helper.addCheck(step12A, v.isInRange(V))
						helper.addCheck(step12B1, v.isValidResidue(a))
						helper.addCheck(step12B2, v.isValidResidue(b))
						helper.addCheck(step12C, c.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, a, b, m)))
						helper.addCheck(step12D, v.powP(v.constants.G, &V).Compare(v.mulP(&a, v.powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
						helper.addCheck(step12E, v.powP(&alpha, &V).Compare(v.mulP(&b, v.powP(&m, &c))))
					}
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
