package core

import (
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateCorrectnessOfDecryptionContestData(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 11, "Correctness of decryptions of contest data (Should not be implemented)")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Should not be implemented for spec 1.91

	v.helpers[helper.VerificationStep] = helper
}
