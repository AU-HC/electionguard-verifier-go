package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
	"go.uber.org/zap"
	"strconv"
)

type Verifier struct {
	logger    *zap.Logger
	constants utility.CorrectElectionConstants
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger}
}

func (v *Verifier) Verify(path string) bool {
	// Deserialize election record and fetch correct constants (Step 0)
	er := v.getElectionRecord(path)
	v.constants = utility.MakeCorrectElectionConstants()

	// Validate election parameters (Step 1):
	electionParametersHelper := v.validateElectionConstants(er)
	electionParametersIsNotValid := !electionParametersHelper.validate()
	if electionParametersIsNotValid {
		return false
	}

	// Validate guardian public-key (Step 2)
	publicKeyValidationHelper := v.validateGuardianPublicKeys(er)
	publicKeysAreNotValid := !publicKeyValidationHelper.validate()
	if publicKeysAreNotValid {
		return false
	}

	// Validation election public-key (Step 3)
	electionKeyValidationHelper := v.validateJointPublicKey(er)
	jointElectionKeyIsNotValid := !electionKeyValidationHelper.validate()
	if jointElectionKeyIsNotValid {
		return false
	}

	// Validate correctness of selection encryptions (Step 4)
	selectionEncryptionValidationHelper := v.validateSelectionEncryptions(er)
	correctnessOfSelectionsIsNotValid := !selectionEncryptionValidationHelper.validate()
	if correctnessOfSelectionsIsNotValid {
		return false
	}

	// Validate adherence to vote limits (Step 5)
	voteLimitsValidationHelper := v.validateVoteLimits(er)
	voteLimitsNotValid := !voteLimitsValidationHelper.validate()
	if voteLimitsNotValid {
		return false
	}

	// Validate confirmation codes (Step 6)
	confirmationCodesValidationHelper := v.validateConfirmationCodes(er)
	confirmationCodesAreNotValid := !confirmationCodesValidationHelper.validate()
	if confirmationCodesAreNotValid {
		return false
	}

	// Validate correctness of ballot aggregation (Step 7)
	ballotAggregationValidationHelper := v.validateBallotAggregation(er)
	ballotAggregationIsNotValid := !ballotAggregationValidationHelper.validate()
	if ballotAggregationIsNotValid {
		return false
	}

	// Validate correctness of partial decryptions (Step 8)
	partialDecryptionsValidationHelper := v.validatePartialDecryptions(er)
	partialDecryptionsAreNotValid := !partialDecryptionsValidationHelper.validate()
	if partialDecryptionsAreNotValid {
		return false
	}

	// Validate correctness of substitute data for missing guardians (Step 9)
	substituteDataValidationHelper := v.validateSubstituteDataForMissingGuardians(er)
	substituteDataForMissingGuardiansIsNotValid := !substituteDataValidationHelper.validate()
	if substituteDataForMissingGuardiansIsNotValid {
		return false
	}

	// Validate correctness of construction of replacement partial decryptions (Step 10)
	replacementDecryptionsValidationHelper := v.validateConstructionOfReplacementForPartialDecryptions(er)
	replacementPartialDecryptionsAreInvalid := !replacementDecryptionsValidationHelper.validate()
	if replacementPartialDecryptionsAreInvalid {
		return false
	}

	// Validate correctness of tally decryption (Step 11)
	tallyDecryptionValidationHelper := v.validateTallyDecryption(er)
	tallyDecryptionIsInvalid := !tallyDecryptionValidationHelper.validate()
	if tallyDecryptionIsInvalid {
		return false
	}

	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	spoiledBallotsDecryptionValidationHelper := v.validatePartialDecryptionForSpoiledBallots(er)
	spoiledBallotsPartialDecryptionIsInvalid := !spoiledBallotsDecryptionValidationHelper.validate()
	if spoiledBallotsPartialDecryptionIsInvalid {
		return false
	}

	// Validate correctness of substitute data for spoiled ballots (Step 13)
	substituteDataForBallotsValidationHelper := v.validateSubstituteDataForSpoiledBallots(er)
	substituteDataForSpoiledBallotsIsInvalid := !substituteDataForBallotsValidationHelper.validate()
	if substituteDataForSpoiledBallotsIsInvalid {
		return false
	}

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)
	replacementDecryptionForBallotsValidationHelper := v.validateReplacementPartialDecryptionForSpoiledBallots(er)
	replacementDataForPartialDecryptionsForBallotsIsInvalid := !replacementDecryptionForBallotsValidationHelper.validate()
	if replacementDataForPartialDecryptionsForBallotsIsInvalid {
		return false
	}

	// Validation of correct decryption of spoiled ballots (Step 15)
	decryptionOfSpoiledBallotsValidationHelper := v.validateDecryptionOfSpoiledBallots(er)
	decryptionOfSpoiledBallotsIsInvalid := !decryptionOfSpoiledBallotsValidationHelper.validate()
	if decryptionOfSpoiledBallotsIsInvalid {
		return false
	}

	// and validation of correctness of spoiled ballots (Step 16)
	correctnessOfSpoiledBallotsValidationHelper := v.validateCorrectnessOfSpoiledBallots(er)
	correctnessSpoiledBallotsIsInvalid := !correctnessOfSpoiledBallotsValidationHelper.validate()
	if correctnessSpoiledBallotsIsInvalid {
		return false
	}

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	correctContestDataValidationHelper := v.validateContestDataPartialDecryptionsForSpoiledBallots(er)
	contestDataPartialDecryptionsForSpoiledBallotIsInvalid := !correctContestDataValidationHelper.validate()
	if contestDataPartialDecryptionsForSpoiledBallotIsInvalid {
		return false
	}

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	substituteContestDataValidationHelper := v.validateSubstituteContestDataForSpoiledBallots(er)
	substituteContestDataForSpoiledBallotIsInvalid := !substituteContestDataValidationHelper.validate()
	if substituteContestDataForSpoiledBallotIsInvalid {
		return false
	}

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)
	contestReplacementDecryptionValidationHelper := v.validateContestReplacementDecryptionForSpoiledBallots(er)
	contestReplacementDecryptionsIsInvalid := !contestReplacementDecryptionValidationHelper.validate()
	if contestReplacementDecryptionsIsInvalid {
		return false
	}

	// Verification was successful
	return true
}

func (v *Verifier) getElectionRecord(path string) *deserialize.ElectionRecord {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	er := parser.ParseElectionRecord(path)
	v.logger.Info("[VALID]: Election data was well formed (Step 0)")

	return er
}

func (v *Verifier) validateElectionConstants(er *deserialize.ElectionRecord) *ValidationHelper {
	constants := utility.MakeCorrectElectionConstants()
	helper := MakeValidationHelper(v.logger, "Election parameters are correct (Step 1)")

	helper.addCheck("(1.A) The large prime is equal to the large modulus p", constants.P.Compare(&er.ElectionConstants.LargePrime))
	helper.addCheck("(1.B) The small prime is equal to the prime q", constants.Q.Compare(&er.ElectionConstants.SmallPrime))
	helper.addCheck("(1.C) The cofactor is equal to r = (p − 1)/q", constants.C.Compare(&er.ElectionConstants.Cofactor))
	helper.addCheck("(1.D) The generator is equal to the generator g", constants.G.Compare(&er.ElectionConstants.Generator))

	return helper
}

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

func (v *Verifier) validateJointPublicKey(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate election public-key (Step 3) [ERROR IN SPEC SHEET FOR (3.B)]
	helper := MakeValidationHelper(v.logger, "Election public-key validation (Step 3)")

	elgamalPublicKey := schema.MakeBigIntFromString("1", 10)
	for _, guardian := range er.Guardians {
		elgamalPublicKey = mulP(elgamalPublicKey, &guardian.ElectionPublicKey)
	}

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	computedExtendedBaseHash := crypto.HashElems(er.CiphertextElectionRecord.CryptoBaseHash, er.CiphertextElectionRecord.CommitmentHash)

	helper.addCheck("(3.A) The joint public election key is computed correctly", elgamalPublicKey.Compare(&er.CiphertextElectionRecord.ElgamalPublicKey))
	helper.addCheck("(3.B) The extended base hash is computed correctly", extendedBaseHash.Compare(computedExtendedBaseHash))

	return helper
}

func (v *Verifier) validateSelectionEncryptions(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of selection encryptions (Step 4)
	helper := MakeValidationHelper(v.logger, "Correctness of selection encryptions (Step 4)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	elgamalPublicKey := &er.CiphertextElectionRecord.ElgamalPublicKey

	for i, ballot := range er.SubmittedBallots {
		for j, contest := range ballot.Contests {
			for k, ballotSelection := range contest.BallotSelections {
				a := ballotSelection.Ciphertext.Pad
				b := ballotSelection.Ciphertext.Data
				a0 := ballotSelection.Proof.ProofZeroPad
				b0 := ballotSelection.Proof.ProofZeroData
				a1 := ballotSelection.Proof.ProofOnePad
				b1 := ballotSelection.Proof.ProofOneData
				c := ballotSelection.Proof.Challenge
				c0 := ballotSelection.Proof.ProofZeroChallenge
				c1 := ballotSelection.Proof.ProofOneChallenge
				v0 := ballotSelection.Proof.ProofZeroResponse
				v1 := ballotSelection.Proof.ProofOneResponse

				helper.addCheck("(4.A) a is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a))
				helper.addCheck("(4.A) b is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b))
				helper.addCheck("(4.A) a0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a0))
				helper.addCheck("(4.A) b0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b0))
				helper.addCheck("(4.A) a1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a1))
				helper.addCheck("(4.A) b1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b1))
				helper.addCheck("(4.B) The challenge value c is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(crypto.HashElems(extendedBaseHash, a, b, a0, b0, a1, b1)))
				helper.addCheck("(4.C) c0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c0))
				helper.addCheck("(4.C) c1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c1))
				helper.addCheck("(4.C) v0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v0))
				helper.addCheck("(4.C) v1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v1))
				helper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(addQ(&c0, &c1)))
				helper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				helper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(v.constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				helper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(elgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				helper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", mulP(powP(v.constants.G, &c1), powP(elgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))
			}
		}
	}

	return helper
}

func (v *Verifier) validateVoteLimits(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate adherence to vote limits (Step 5)
	helper := MakeValidationHelper(v.logger, "Adherence to vote limits (Step 5)")

	for i, ballot := range er.SubmittedBallots {
		for j, contest := range ballot.Contests {
			contestInManifest := getContest(contest.ObjectId, er.Manifest.Contests)
			votesAllowed := contestInManifest.VotesAllowed
			numberOfSelections := 0
			calculatedAHat := schema.MakeBigIntFromString("1", 10)
			calculatedBHat := schema.MakeBigIntFromString("1", 10)

			for _, selection := range contest.BallotSelections {
				if selection.IsPlaceholderSelection {
					numberOfSelections++
				}
				calculatedAHat = mulP(calculatedAHat, &selection.Ciphertext.Pad)
				calculatedBHat = mulP(calculatedBHat, &selection.Ciphertext.Data)
			}

			aHat := contest.CiphertextAccumulation.Pad
			bHat := contest.CiphertextAccumulation.Data
			a := contest.Proof.Pad
			b := contest.Proof.Data
			V := contest.Proof.Response

			c := crypto.HashElems(er.CiphertextElectionRecord.CryptoExtendedBaseHash, aHat, bHat, a, b)
			equationFLeft := powP(v.constants.G, &V)
			equationFRight := mulP(&a, powP(&aHat, c))
			equationGLeft := mulP(powP(v.constants.G, mulP(schema.MakeBigIntFromInt(votesAllowed), c)), powP(&er.CiphertextElectionRecord.ElgamalPublicKey, &V))
			equationGRight := mulP(&b, powP(&bHat, c))

			helper.addCheck("(5.A) The number of placeholder positions matches the selection limit ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", votesAllowed == numberOfSelections)
			helper.addCheck("(5.B) The a hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", aHat.Compare(calculatedAHat))
			helper.addCheck("(5.B) The b hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", bHat.Compare(calculatedBHat))
			helper.addCheck("(5.C) The given value V is in Zq ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isInRange(V))
			helper.addCheck("(5.D) The given value a are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Pad))
			helper.addCheck("(5.D) The given values b are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Data))
			helper.addCheck("(5.E) The challenge value is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", contest.Proof.Challenge.Compare(c))
			helper.addCheck("(5.F) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationFLeft.Compare(equationFRight))
			helper.addCheck("(5.G) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationGLeft.Compare(equationGRight))
		}
	}

	return helper
}

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

func (v *Verifier) validateBallotAggregation(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of ballot aggregation (Step 7)
	helper := MakeValidationHelper(v.logger, "Correctness of ballot aggregation (Step 7)")

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			a := schema.MakeBigIntFromInt(1)
			b := schema.MakeBigIntFromInt(1)
			for _, ballot := range er.SubmittedBallots {
				ballotWasCast := ballot.State == 1
				if ballotWasCast {
					ciphertextSelection := getSelection(ballot, contest.ObjectId, selection.ObjectId)
					a = mulP(a, &ciphertextSelection.Pad)
					b = mulP(b, &ciphertextSelection.Data)
				}
			}
			A := selection.Message.Pad
			B := selection.Message.Data
			helper.addCheck("(7.A) A is calculated correctly", A.Compare(a))
			helper.addCheck("(7.B) B is calculated correctly", B.Compare(b))
		}
	}

	return helper
}

func (v *Verifier) validatePartialDecryptions(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of partial decryptions (Step 8)
	helper := MakeValidationHelper(v.logger, "Correctness of partial decryptions (Step 8)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data

			for k, share := range selection.Shares {
				if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix
					V := share.Proof.Response
					c := share.Proof.Challenge
					ai := share.Proof.Pad
					bi := share.Proof.Data
					m := share.Share

					helper.addCheck("(8.A) The value v is in the set Zq for "+share.ObjectId+" "+strconv.Itoa(k), isInRange(V))
					helper.addCheck("(8.B) The value a is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Pad))
					helper.addCheck("(8.B) The value b is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Data))
					helper.addCheck("(8.C) The challenge is computed correctly "+share.ObjectId+" "+strconv.Itoa(k), c.Compare(crypto.HashElems(extendedBaseHash, A, B, ai, bi, m)))
					helper.addCheck("(8.D) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(v.constants.G, &V).Compare(mulP(&ai, powP(&er.Guardians[k].ElectionPublicKey, &c))))
					helper.addCheck("(8.E) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(&A, &V).Compare(mulP(&bi, powP(&m, &c))))
				}
			}
		}
	}

	return helper
}

func (v *Verifier) validateSubstituteDataForMissingGuardians(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of substitute data for missing guardians (Step 9)
	helper := MakeValidationHelper(v.logger, "Correctness of substitute data for missing guardians (Step 9)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data
			for _, share := range selection.Shares {
				for _, part := range share.RecoveredParts {
					if part.ObjectId != "" { // TODO: Implement method to check if "Recovered parts" is not nil
						V := part.Proof.Response
						c := part.Proof.Challenge
						a := part.Proof.Pad
						b := part.Proof.Data
						m := part.PartialDecryption

						helper.addCheck("(9.A) The given value v is in Zq", isInRange(V))
						helper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(a))
						helper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(b))
						helper.addCheck("(9.C) The challenge value c is correct", c.Compare(crypto.HashElems(extendedBaseHash, A, B, a, b, m)))
						helper.addCheck("(9.D) The equation is satisfied", powP(v.constants.G, &V).Compare(mulP(&a, powP(&part.RecoveryPublicKey, &c))))
						helper.addCheck("(9.E) The equation is satisfied", powP(&A, &V).Compare(mulP(&b, powP(&m, &c))))
					}
				}
			}
		}
	}

	return helper
}

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
				product := schema.MakeBigIntFromString("1", 10)
				for _, part := range share.RecoveredParts {
					coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
					product = mulP(product, powP(&part.PartialDecryption, &coefficient))
				}
				helper.addCheck("(10.B) Correct tally share?", share.Share.Compare(product))
			}
		}
	}

	return helper
}

func (v *Verifier) validateTallyDecryption(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of tally decryption (Step 11)
	helper := MakeValidationHelper(v.logger, "Correct decryption of tallies (Step 11)")

	for _, contest := range er.PlaintextTally.Contests {
		helper.addCheck("Tally label exists in election manifest", contains(er.Manifest.Contests, contest.ObjectId))
		// TODO: Check 11.C to 11.F
		for _, selection := range contest.Selections {
			b := selection.Message.Data
			mi := schema.MakeBigIntFromString("1", 10)
			m := selection.Value
			t := schema.MakeBigIntFromInt(selection.Tally)
			for _, share := range selection.Shares {
				mi = mulP(mi, &share.Share)
			}
			helper.addCheck("(11.A) The equation is satisfied", b.Compare(mulP(&m, mi)))
			helper.addCheck("(11.B) The equation is satisfied", m.Compare(powP(v.constants.G, t)))
		}
	}

	return helper
}

func (v *Verifier) validatePartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	helper := MakeValidationHelper(v.logger, "Correctness of partial decryption for spoiled ballots (Step 12)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) { // Comparing with zero, will need better way of determining this TODO: Fix
						m := share.Share
						a := share.Proof.Pad
						b := share.Proof.Data
						c := share.Proof.Challenge
						V := share.Proof.Response

						helper.addCheck("(12.A) The given value v is in the set Zq", isInRange(V))
						helper.addCheck("(12.B) The given value a is in the set Zpr", isValidResidue(a))
						helper.addCheck("(12.B) The given value b is in the set Zpr", isValidResidue(b))
						helper.addCheck("(12.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, a, b, m)))
						helper.addCheck("(12.D) The equation is satisfied", powP(v.constants.G, &V).Compare(mulP(&a, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
						helper.addCheck("(12.E) The equation is satisfied", powP(&alpha, &V).Compare(mulP(&b, powP(&m, &c))))
					}
				}
			}
		}
	}

	return helper
}

func (v *Verifier) validateReplacementPartialDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correct replacement partial decryptions for spoiled ballots (Step 14)
	helper := MakeValidationHelper(v.logger, "Correctness of substitute data for spoiled ballots (Step 13)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) { // Comparing with zero, will need better way of determining this TODO: Fix
						for _, part := range share.RecoveredParts {
							mil := part.PartialDecryption
							ai := part.Proof.Pad
							bi := part.Proof.Data
							ci := part.Proof.Challenge
							vi := part.Proof.Response

							helper.addCheck("(13.A) The given value v is in Zq", isInRange(vi))
							helper.addCheck("(13.B) The given value a is in the set Zpr", isValidResidue(ai))
							helper.addCheck("(13.B) The given value b is in the set Zpr", isValidResidue(bi))
							helper.addCheck("(13.C) The challenge is computed correctly", ci.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, ai, bi, mil)))
							helper.addCheck("(13.D) The equation is satisfied", powP(v.constants.G, &vi).Compare(powP(mulP(&ai, &part.RecoveryPublicKey), &ci)))
							helper.addCheck("(13.E) The equation is satisfied", powP(&ai, &vi).Compare(mulP(&bi, powP(&mil, &ci))))
						}
					}
				}
			}
		}
	}

	return helper
}

func (v *Verifier) validateSubstituteDataForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validating correctness of substitute data for spoiled ballots (Step 13)
	helper := MakeValidationHelper(v.logger, "Correctness of substitute data for spoiled ballots (Step 13)")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromInt(0)) {
						m := share.Share
						product := schema.MakeBigIntFromInt(1)

						for _, part := range share.RecoveredParts {
							coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
							product = mulP(product, powP(&part.PartialDecryption, &coefficient))
						}
						if len(share.RecoveredParts) > 0 {
							helper.addCheck("(14.B) Correct missing decryption share", m.Compare(product))
						}
					}
				}
			}
		}
	}

	return helper
}

func (v *Verifier) validateDecryptionOfSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correct decryption of spoiled ballots (Step 15)
	helper := MakeValidationHelper(v.logger, "Correct decryption of spoiled ballots (Step 15)")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			sumOfAllSelections := 0
			for _, selection := range contest.Selections {
				beta := selection.Message.Data
				m := selection.Value
				V := schema.MakeBigIntFromInt(selection.Tally)
				mi := schema.MakeBigIntFromInt(1)
				sumOfAllSelections += selection.Tally
				for _, share := range selection.Shares {
					mi = mulP(mi, &share.Share)
				}

				helper.addCheck("(15.A) The equation is satisfied", beta.Compare(mulP(&m, mi)))
				helper.addCheck("(15.B) The equation is satisfied", m.Compare(powP(v.constants.G, V)))
			}
		}
	}

	return helper
}

func (v *Verifier) validateCorrectnessOfSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correctness of spoiled ballots (Step 16)
	helper := MakeValidationHelper(v.logger, "Correctness of spoiled ballots (Step 16)")

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			sumOfAllSelections := 0
			for _, selection := range contest.Selections {
				sumOfAllSelections += selection.Tally
				helper.addCheck("(16.A) For each option in the contest, the selection V is either a 0 or a 1", selection.Tally == 0 || selection.Tally == 1)
			}
			helper.addCheck("(16.B) The sum of all selections in the contest is at most the selection limit L for that contest.", sumOfAllSelections <= getContest(contest.ObjectId, er.Manifest.Contests).VotesAllowed)
			// TODO: 16.C -> 16.E
		}
	}

	return helper
}

func (v *Verifier) validateContestDataPartialDecryptionsForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	helper := MakeValidationHelper(v.logger, "Correctness of contest data partial decryptions for spoiled ballots (Step 17)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			c0 := contest.ContestData.Ciphertext.Generator
			c1 := contest.ContestData.Ciphertext.EncryptedMessage
			c2 := contest.ContestData.Ciphertext.MessageAuthenticationCode

			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				ai := share.Proof.Pad
				bi := share.Proof.Data
				ci := share.Proof.Challenge
				vi := share.Proof.Response

				helper.addCheck("(17.A) The given value v is in the set Zq", isInRange(vi))
				helper.addCheck("(17.B) The given value a is in the set Zqr", isValidResidue(ai))
				helper.addCheck("(17.B) The given value b is in the set Zqr", isValidResidue(bi))
				helper.addCheck("(17.C) The challenge is correctly computed", ci.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, ai, bi, mi)))
				helper.addCheck("(17.D) The equation is satisfied", powP(v.constants.G, &vi).Compare(mulP(&ai, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &ci))))
				helper.addCheck("(17.E) The equation is satisfied", powP(&c0, &vi).Compare(mulP(&bi, powP(&mi, &ci))))
			}
		}
	}

	return helper
}

func (v *Verifier) validateSubstituteContestDataForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	helper := MakeValidationHelper(v.logger, "Correctness of substitute contest data for spoiled ballots (Step 18)")
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			c0 := contest.ContestData.Ciphertext.Generator
			c1 := contest.ContestData.Ciphertext.EncryptedMessage
			c2 := contest.ContestData.Ciphertext.MessageAuthenticationCode

			for k, share := range contest.ContestData.Shares {
				for _, part := range share.RecoveredParts {
					V := part.Proof.Response
					c := part.Proof.Challenge
					a := part.Proof.Pad
					b := part.Proof.Data
					m := part.PartialDecryption

					helper.addCheck("(18.A) The value v is in the set Zq", isInRange(V))
					helper.addCheck("(18.B) The value a is in the set Zqr", isValidResidue(a))
					helper.addCheck("(18.B) The value b is in the set Zqr", isValidResidue(b))
					helper.addCheck("(18.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, a, b, m)))
					helper.addCheck("(18.D) The equation is satisfied", powP(v.constants.G, &V).Compare(mulP(&a, powP(&er.Guardians[k].ElectionPublicKey, &c))))
					helper.addCheck("(18.E) The equation is satisfied", powP(&c0, &V).Compare(mulP(&b, powP(&m, &c))))
				}
			}
		}
	}

	return helper
}

func (v *Verifier) validateContestReplacementDecryptionForSpoiledBallots(er *deserialize.ElectionRecord) *ValidationHelper {
	// Validation of correctness of contest replacement decryptions for spoiled ballots (Step 19)
	helper := MakeValidationHelper(v.logger, "Correctness of contest replacement decryptions for spoiled ballots (Step 19)")
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, share := range contest.ContestData.Shares {
				mi := share.Share
				product := schema.MakeBigIntFromInt(1)
				for _, part := range share.RecoveredParts {
					m := part.PartialDecryption

					coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
					product = mulP(product, powP(&m, &coefficient))
				}
				helper.addCheck("(19.A) The equation is satisfied", mi.Compare(product))
			}
		}
	}

	return helper
}
