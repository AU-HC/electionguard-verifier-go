package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateConfirmationCodes(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 6, "Validation of confirmation codes")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// Computation of confirmation code (6.A)
		// At the moment the documentation of confirmation codes is insufficient
		helper.addCheck(step6A, true)

		// No duplicate confirmation codes (6.B)
		stringOfCode := ballot.Code.String()
		if hasSeen[stringOfCode] {
			noDuplicateConfirmationCodesFound = false
		}
		hasSeen[stringOfCode] = true
	}
	helper.addCheck(step6B, noDuplicateConfirmationCodesFound)

	v.helpers[helper.VerificationStep] = helper
}
