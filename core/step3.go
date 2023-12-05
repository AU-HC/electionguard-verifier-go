package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateElectionPublicKey(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 3, "Election public-key is correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Calculating the joint election public-key and checking that the individual keys are valid
	electionPublicKey := schema.MakeBigIntFromString("1", 10)
	one := schema.MakeBigIntFromString("1", 10)
	for _, guardian := range er.Guardians {
		electionPublicKey = v.mulP(electionPublicKey, &guardian.Key)
		errorString := "(GuardianID:" + guardian.ObjectID + ")"
		helper.addCheck("(3.A) The guardian public-key is not in Z^r_p.", v.isValidResidue(guardian.Key) && !one.Compare(&guardian.Key), errorString)
	}

	// Validating the joint election public-key
	helper.addCheck("(3.B) The election public key is not correct.", er.CiphertextElectionRecord.ElgamalPublicKey.Compare(electionPublicKey))

	v.helpers[helper.VerificationStep] = helper
}
