package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"strconv"
)

func (v *Verifier) validateGuardianPublicKeys(er *deserialize.ElectionRecord) *ValidationHelper {
	helper := MakeValidationHelper(v.logger, "Guardian public-key validation (Step 2)")

	for i, guardian := range er.Guardians {
		for j, proof := range guardian.ElectionProofs {
			// (2.A)
			hash := crypto.HashElems(guardian.ElectionCommitments[j], proof.Commitment)
			helper.addCheck("(2.A) The challenge is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", proof.Challenge.Compare(hash))

			// (2.B)
			left := powP(v.constants.G, &proof.Response)
			right := mulP(powP(&guardian.ElectionCommitments[j], &proof.Challenge), &proof.Commitment)
			helper.addCheck("(2.B) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", left.Compare(right))
		}
	}

	return helper
}
