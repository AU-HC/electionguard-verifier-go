package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateConstructionOfReplacementForPartialDecryptions(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 10, "Correctness of construction of replacement partial decryptions")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// 10.A, 14.A? TODO: Refactor
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
		helper.addCheck(step10A, productJ.Compare(productJMinusL))
	}

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			for _, share := range selection.Shares {
				if len(share.RecoveredParts) > 0 {

					product := schema.MakeBigIntFromInt(1)

					for _, part := range share.RecoveredParts {
						coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
						product = v.mulP(product, v.powP(&part.Share, &coefficient))
					}

					helper.addCheck(step10B, share.Share.Compare(product))
				}
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}
