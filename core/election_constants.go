package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
	"time"
)

func (v *Verifier) validateElectionConstants(er *deserialize.ElectionRecord) {
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 1, "Election parameters are correct")
	start := time.Now()

	constants := utility.MakeCorrectElectionConstants()
	helper.addCheck(step1A, constants.P.Compare(&er.ElectionConstants.LargePrime))
	helper.addCheck(step1B, constants.Q.Compare(&er.ElectionConstants.SmallPrime))
	helper.addCheck(step1C, constants.C.Compare(&er.ElectionConstants.Cofactor))
	helper.addCheck(step1D, constants.G.Compare(&er.ElectionConstants.Generator))

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 1 took: " + time.Since(start).String())
}
