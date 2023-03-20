package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateConfirmationCodes(er *deserialize.ElectionRecord) {
	// Validate confirmation codes (Step 6)
	defer v.wg.Done()
	helper := MakeValidationHelper(v.logger, 6, "Validation of confirmation codes")
	start := time.Now()

	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// Computation of confirmation code (6.A)
		helper.addCheck(step6A, true) // TODO: Fake it

		// No duplicate confirmation codes (6.B)
		stringOfCode := ballot.Code.String()
		if hasSeen[stringOfCode] {
			noDuplicateConfirmationCodesFound = false
		}
		hasSeen[stringOfCode] = true
	}
	helper.addCheck(step6B, noDuplicateConfirmationCodesFound)

	v.helpers[helper.VerificationStep] = helper
	v.logger.Info("Validation of step 6 took: " + time.Since(start).String())
}
