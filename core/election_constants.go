package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
)

func (v *Verifier) validateElectionConstants(er *deserialize.ElectionRecord) {
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 1, "Election parameters are correct")

	constants := utility.MakeCorrectElectionConstants()
	helper.addCheck("(1.A) The large prime is equal to the large modulus p", constants.P.Compare(&er.ElectionConstants.LargePrime))
	helper.addCheck("(1.B) The small prime is equal to the prime q", constants.Q.Compare(&er.ElectionConstants.SmallPrime))
	helper.addCheck("(1.C) The cofactor is equal to r = (p − 1)/q", constants.C.Compare(&er.ElectionConstants.Cofactor))
	helper.addCheck("(1.D) The generator is equal to the generator g", constants.G.Compare(&er.ElectionConstants.Generator))

	v.helpers[helper.VerificationStep] = helper
}
