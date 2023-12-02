package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateAdherenceToVoteLimits(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 6, "Adherence to vote limits")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Should not be implemented for spec 1.91

	v.helpers[helper.VerificationStep] = helper
}
