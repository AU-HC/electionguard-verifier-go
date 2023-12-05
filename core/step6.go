package core

import (
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateAdherenceToVoteLimits(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 6, "Adherence to vote limits (Should not be implemented)")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Should not be implemented for spec 1.91

	v.helpers[helper.VerificationStep] = helper
}
