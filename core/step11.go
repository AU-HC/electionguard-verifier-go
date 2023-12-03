package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateCorrectnessOfDecryptionContestData(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 11, "Correctness of decryptions of contest data")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Should not be implemented for spec 1.91

	v.helpers[helper.VerificationStep] = helper
}
