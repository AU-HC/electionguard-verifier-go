package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateElectionConstants(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 1, "Election parameters are correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	helper.addCheck(step1A, v.constants.P.Compare(&er.ElectionConstants.LargePrime))
	helper.addCheck(step1B, v.constants.Q.Compare(&er.ElectionConstants.SmallPrime))
	helper.addCheck(step1C, v.constants.C.Compare(&er.ElectionConstants.Cofactor))
	helper.addCheck(step1D, v.constants.G.Compare(&er.ElectionConstants.Generator))

	v.helpers[helper.VerificationStep] = helper
}
