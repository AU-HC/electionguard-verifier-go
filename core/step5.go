package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"fmt"
	"time"
)

func (v *Verifier) validateSelectionEncryptions(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 5, "Selection encryptions are correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.MockBallots
	chunkSize := v.verifierStrategy.getBallotChunkSize(len(ballots))

	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		helper.wg.Add(1)
		go v.validateSelectionEncryptionForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSelectionEncryptionForSlice(helper *ValidationHelper, ballots []schema.MockBallot, er *schema.ElectionRecord) {
	defer helper.wg.Done()

	q := er.ElectionConstants.SmallPrime
	g := &er.ElectionConstants.Generator
	k := schema.MakeBigIntFromString("FE211456D2A9E67D28009C885B9052B4999F2F97A930C2557AE9DD346BDD20F8FEB09FD8CD693D7FAC646A7683D028B43A87C224E62DC76832272223582C768CEDA25A17A6ADB0103CCD3B18175E0B226D54DA939C651497E091D1FEAB34EC45490E287E8A14C06397CEDD0FFA3C2B410ACD443E1961FCE0135ED1231499242B2E2A53A9618BB7C82BF76A0C035281A1C6E7ABCEEAADFDDC53969F350057C8723B96D00B3F9C72FF07AAE0BF94ED52D160812BDE79F131BD6708B30C2934E69E8C085F688B0692E5DCD908F984845DCDEAA33355C1B5D0283FFAF727D828D5D2DC7500A9342C3E0B85B7CD78BFA21B7D31A8FB4EF835DBD42BDFB9B423BDA11FCB552BC0B1FBD62CFD91444FB176FB9DB9B306A1ECC7944DD2E2CC269DA3C0A87E7CDB864838B5C385B04D6D6413505495F46A85C8C8AB04271186565C671B4050AF14871E4B996CC3261806F30C34C40222956B17830B299A1C988203BEF047B9585072281075F850B591BDB345FCE8C86FCFA68C2B5AFDFD4D170DB867466D722139A13D017DD783BC66050952951DC9002CB53C8ED5979E0794A6879076061A86E4248C7EF3AD91C50F4CA17ABB52E7AF6A376A34AB4AF9BF685B15035C8966F8B903FCCAB1B9D5BC5AB4AF0BE3EC4DDC822C7B683E61BBD1A84275398283356377998F2C7EF5A53C6133EF0A2C138AB0F9C1052C50F565F471C118F234FE", 16)
	extendedBaseHash := schema.MakeBigIntFromString("723C8D7FA2AD8015C74A9AFA99D8AF60B6A6E6841C0433547187852E63A0DDF9", 16)

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				alpha := selection.Ciphertext.Pad
				beta := selection.Ciphertext.Data

				fmt.Println(alpha.Text(16))

				toBeHashed := []interface{}{*k, selection.Ciphertext.Pad, selection.Ciphertext.Data}
				computedC := schema.MakeBigIntFromString("0", 10)
				for j, proof := range selection.Proof.Proofs {
					cj := proof.Challenge
					computedC = v.addQ(computedC, &cj)

					vj := proof.Response
					wj := v.subQ(&vj, v.mulQ(schema.MakeBigIntFromString(fmt.Sprintf("%d", j), 10), &cj))

					aj := v.mulP(v.powP(g, &vj), v.powP(&alpha, &cj))
					bj := v.mulP(v.powP(k, wj), v.powP(&beta, &cj))

					toBeHashed = append(toBeHashed, *aj)
					toBeHashed = append(toBeHashed, *bj)

					helper.addCheck("(5.B) ...", v.isInRange(cj))
					helper.addCheck("(5.C) ...", v.isInRange(vj))
				}
				c := crypto.HMAC(q, *extendedBaseHash, 0x21, toBeHashed...)

				helper.addCheck("(5.A) alpha is in Z_p^r", v.isValidResidue(alpha))
				helper.addCheck("(5.A) beta is in Z_p^r", v.isValidResidue(beta))
				helper.addCheck("(5.D) challenge is computed correctly", computedC.Compare(c))
			}
		}
	}
}
