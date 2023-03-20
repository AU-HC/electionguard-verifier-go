package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateSubstituteDataForSpoiledBallots(er *deserialize.ElectionRecord) {
	// Validating correctness of substitute data for spoiled ballots (Step 13)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 13, "Correctness of substitute data for spoiled ballots")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	start := time.Now()

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data

				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) {

						for _, part := range share.RecoveredParts {
							mil := part.Share
							ai := part.Proof.Pad
							bi := part.Proof.Data
							ci := part.Proof.Challenge
							vi := part.Proof.Response

							helper.addCheck(step13A, v.isInRange(vi))
							helper.addCheck(step13B1, v.isValidResidue(ai))
							helper.addCheck(step13B2, v.isValidResidue(bi))
							helper.addCheck(step13C, ci.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, ai, bi, mil)))
							helper.addCheck(step13D, v.powP(v.constants.G, &vi).Compare(v.powP(v.mulP(&ai, &part.RecoveryPublicKey), &ci)))
							helper.addCheck(step13E, v.powP(&ai, &vi).Compare(v.mulP(&bi, v.powP(&mil, &ci))))

						}
					}
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 13 took: " + time.Since(start).String())
}
