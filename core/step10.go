package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateCorrectnessOfTallyDecryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 10, "Correct decryption of tallies")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Should not be implemented for spec 1.91

	v.helpers[helper.VerificationStep] = helper
}
