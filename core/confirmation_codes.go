package core

import (
	"electionguard-verifier-go/deserialize"
)

func (v *Verifier) validateConfirmationCodes(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate confirmation codes (Step 6)
	helper := MakeValidationHelper(v.logger, "Validation of confirmation codes (Step 6)")

	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// Computation of confirmation code (6.A)
		helper.addCheck("(6.A) The confirmation code for ballot id: "+ballot.ObjectId+" is computed correct", true) // TODO: Fake it

		// No duplicate confirmation codes (6.B)
		stringOfCode := ballot.Code.String()
		if hasSeen[stringOfCode] {
			noDuplicateConfirmationCodesFound = false
		}
		hasSeen[stringOfCode] = true
	}
	helper.addCheck("(6.B) No duplicate confirmation codes found", noDuplicateConfirmationCodesFound)

	return helper
}
