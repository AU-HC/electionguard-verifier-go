package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateConfirmationCodes(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 7, "No duplicate confirmation codes found")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Create map of seen confirmation codes (could also create a set by map[string]interface{}, however this would only save 1bit per code)
	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// No duplicate confirmation codes (6.B)
		stringOfCode := ballot.Code.String()
		if hasSeen[stringOfCode] {
			noDuplicateConfirmationCodesFound = false
		}
		hasSeen[stringOfCode] = true
	}
	helper.addCheck("(7.C) Duplicate confirmation codes were found", noDuplicateConfirmationCodesFound)

	v.helpers[helper.VerificationStep] = helper
}
