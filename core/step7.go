package core

import (
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateConfirmationCodes(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 7, "No duplicate confirmation codes found")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Create map of seen confirmation codes
	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// No duplicate confirmation codes (6.B)
		if hasSeen[ballot.Code] {
			noDuplicateConfirmationCodesFound = false
			helper.addCheck("(7.C) Duplicate confirmation codes were found.", noDuplicateConfirmationCodesFound, ballot.Code)
		}
		hasSeen[ballot.Code] = true
	}

	v.helpers[helper.VerificationStep] = helper
}
