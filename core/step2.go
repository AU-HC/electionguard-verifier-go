package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"fmt"
	"time"
)

func (v *Verifier) validateGuardianPublicKeys(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 2, "Guardian public-key is correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	g := &er.ElectionConstants.Generator
	q := &er.ElectionConstants.SmallPrime

	for _, guardian := range er.Guardians {
		for i, schnorrProof := range guardian.CoefficientsProofs {
			// Computing validity of Schnorr proof
			response := schnorrProof.Response
			leftSchnorr := v.powP(g, &response)
			rightSchnorr := v.mulP(&schnorrProof.Commitment, v.powP(&schnorrProof.PublicKey, &schnorrProof.Challenge))

			// Computing challenge
			expectedC := schnorrProof.Challenge
			calculatedC := crypto.Hash(q, schnorrProof.PublicKey, schnorrProof.Commitment)

			// Adding checks
			errorString := fmt.Sprintf("(GuardianID:"+guardian.ObjectID+", Coefficient:%d)", i)
			helper.addCheck("(2.1) The Schnorr proof is not valid.", leftSchnorr.Compare(rightSchnorr), errorString)
			helper.addCheck("(2.A) The value K_{i, j} is not in Z^r_p.", v.isValidResidue(schnorrProof.PublicKey), errorString)
			helper.addCheck("(2.B) The value v_{i, j} is not in Z_q.", v.isInRange(response), errorString)
			helper.addCheck("(2.C) The challenge is not computed correctly", expectedC.Compare(calculatedC), errorString)
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
