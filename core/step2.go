package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateGuardianPublicKeys(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 2, "Guardian public-key is correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	g := &er.ElectionConstants.Generator
	q := &er.ElectionConstants.SmallPrime

	for _, guardian := range er.Guardians {
		for _, schnorrProof := range guardian.CoefficientsProofs {
			// Computing validity of Schnorr proof
			response := schnorrProof.Response
			leftSchnorr := v.powP(g, &response)
			rightSchnorr := v.mulP(&schnorrProof.Commitment, v.powP(&schnorrProof.PublicKey, &schnorrProof.Challenge))

			// Computing challenge
			expectedC := schnorrProof.Challenge
			calculatedC := crypto.Hash1(q, schnorrProof.PublicKey, schnorrProof.Commitment)

			// Adding checks
			helper.addCheck("(2.1) The Schnorr proof is not valid.", leftSchnorr.Compare(rightSchnorr))
			helper.addCheck("(2.A) The value K_{i, j} is not in Z^r_p.", v.isValidResidue(schnorrProof.PublicKey))
			helper.addCheck("(2.B) The value v_{i, j} is not in Z_q.", v.isInRange(response))
			helper.addCheck("(2.C) The challenge is not computed correctly", expectedC.Compare(calculatedC))
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
