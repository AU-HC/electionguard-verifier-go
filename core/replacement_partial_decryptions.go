package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func (v *Verifier) validateConstructionOfReplacementForPartialDecryptions(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of construction of replacement partial decryptions (Step 10)
	helper := MakeValidationHelper(v.logger, "Correctness of construction of replacement partial decryptions (Step 10)")

	// 10.A TODO: Refactor
	for l, wl := range er.CoefficientsValidationSet.Coefficients {
		productJ := schema.MakeBigIntFromInt(1)
		productJMinusL := schema.MakeBigIntFromInt(1)

		for j := range er.CoefficientsValidationSet.Coefficients {
			if j != l {
				jInt := schema.MakeBigIntFromString(j, 10)
				lInt := schema.MakeBigIntFromString(l, 10)
				productJ = mul(productJ, jInt)
				productJMinusL = mul(productJMinusL, sub(jInt, lInt))
			}
		}
		productJ = modQ(productJ)
		productJMinusL = modQ(mul(&wl, productJMinusL))
		helper.addCheck("(10.A) Coefficient check for guardian "+l, productJ.Compare(productJMinusL))
	}

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			for _, share := range selection.Shares {
				if share.Proof.Usage == "" { // TODO: Need better way of checking for nil
					product := schema.MakeBigIntFromInt(1)

					for _, part := range share.RecoveredParts {
						coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
						product = mulP(product, powP(&part.Share, &coefficient))
					}

					helper.addCheck("(10.B) Correct tally share ", share.Share.Compare(product))
				}
			}
		}
	}

	return helper
}
