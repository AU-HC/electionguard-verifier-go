package core

import (
	"electionguard-verifier-go/deserialize"
	"time"
)

func (v *Verifier) validateConfirmationCodes(er *deserialize.ElectionRecord) {
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
		}
		hasSeen[ballot.Code] = true
	}

	// TODO: send the one that is duplicate?
	helper.addCheck("(7.C) Duplicate confirmation codes were found", noDuplicateConfirmationCodesFound)

	v.helpers[helper.VerificationStep] = helper
}
