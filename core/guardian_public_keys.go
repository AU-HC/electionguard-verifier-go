package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateGuardianPublicKeys(er *deserialize.ElectionRecord) {
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 2, "Guardian public-key validation")
	start := time.Now()

	// TODO: Check j
	for _, guardian := range er.Guardians {
		for j, proof := range guardian.ElectionProofs {
			// (2.A)
			hash := crypto.HashElems(guardian.ElectionCommitments[j], proof.Commitment)
			helper.addCheck(step2A, proof.Challenge.Compare(hash))

			// (2.B)
			left := v.powP(v.constants.G, &proof.Response)
			right := v.mulP(v.powP(&guardian.ElectionCommitments[j], &proof.Challenge), &proof.Commitment)
			helper.addCheck(step2B, left.Compare(right))
		}
	}

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 2 took: " + time.Since(start).String())
}
