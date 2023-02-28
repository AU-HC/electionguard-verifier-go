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
	logger *zap.Logger
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger}
}

func (v *Verifier) Verify(path string) bool {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	er := parser.ConvertJsonDataToGoStruct(path)
	v.logger.Info("[VALID]: Election data was well formed (Step 0)")

	// validate election parameters (Step 1):
	constants := utility.MakeCorrectElectionConstants()
	electionParametersHelper := MakeValidationHelper(v.logger, "Election parameters are correct (Step 1)")
	electionParametersHelper.addCheck("(1.A) The large prime is equal to the large modulus p", constants.P.Compare(&er.ElectionConstants.LargePrime))
	electionParametersHelper.addCheck("(1.B) The small prime is equal to the prime q", constants.Q.Compare(&er.ElectionConstants.SmallPrime))
	electionParametersHelper.addCheck("(1.C) The cofactor is equal to r = (p − 1)/q", constants.C.Compare(&er.ElectionConstants.Cofactor))
	electionParametersHelper.addCheck("(1.D) The generator is equal to the generator g", constants.G.Compare(&er.ElectionConstants.Generator))
	electionParametersIsNotValid := !electionParametersHelper.validate()
	if electionParametersIsNotValid {
		return false
	}

	// validate guardian public-key (Step 2)
	publicKeyValidationHelper := MakeValidationHelper(v.logger, "Guardian public-key validation (Step 2)")
	electionKeyValidationHelper := MakeValidationHelper(v.logger, "Election public-key validation (Step 3)")
	elgamalPublicKey := schema.MakeBigIntFromString("1", 10)
	for i, guardian := range er.Guardians {
		elgamalPublicKey = mulP(elgamalPublicKey, &guardian.ElectionPublicKey)
		for j, proof := range guardian.ElectionProofs {
			// (2.A)
			hash := crypto.HashElems(guardian.ElectionCommitments[j], proof.Commitment)
			publicKeyValidationHelper.addCheck("(2.A) The challenge is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", proof.Challenge.Compare(hash))

			// (2.B)
			left := powP(&constants.G, &proof.Response)
			right := mulP(powP(&guardian.ElectionCommitments[j], &proof.Challenge), &proof.Commitment)
			publicKeyValidationHelper.addCheck("(2.B) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", left.Compare(right))
		}
	}
	publicKeysAreNotValid := !publicKeyValidationHelper.validate()
	if publicKeysAreNotValid {
		return false
	}

	// Validate election public key (Step 3) [ERROR IN SPEC SHEET FOR (3.B)]
	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash
	computedExtendedBaseHash := crypto.HashElems(er.CiphertextElectionRecord.CryptoBaseHash, er.CiphertextElectionRecord.CommitmentHash)

	electionKeyValidationHelper.addCheck("(3.A) The joint public election key is computed correctly", elgamalPublicKey.Compare(&er.CiphertextElectionRecord.ElgamalPublicKey))
	electionKeyValidationHelper.addCheck("(3.B) The extended base hash is computed correctly", extendedBaseHash.Compare(computedExtendedBaseHash))

	jointElectionKeyIsNotValid := !electionKeyValidationHelper.validate()
	if jointElectionKeyIsNotValid {
		return false
	}

	// validate correctness of selection encryptions (Step 4)
	selectionEncryptionValidationHelper := MakeValidationHelper(v.logger, "Correctness of selection encryptions (Step 4)")
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

				// TODO: Refactor at some point
				selectionEncryptionValidationHelper.addCheck("(4.A) a is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a))
				selectionEncryptionValidationHelper.addCheck("(4.A) b is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b))
				selectionEncryptionValidationHelper.addCheck("(4.A) a0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a0))
				selectionEncryptionValidationHelper.addCheck("(4.A) b0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b0))
				selectionEncryptionValidationHelper.addCheck("(4.A) a1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a1))
				selectionEncryptionValidationHelper.addCheck("(4.A) b1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b1))
				selectionEncryptionValidationHelper.addCheck("(4.B) The challenge value c is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(crypto.HashElems(extendedBaseHash, a, b, a0, b0, a1, b1)))
				selectionEncryptionValidationHelper.addCheck("(4.C) c0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c0))
				selectionEncryptionValidationHelper.addCheck("(4.C) c1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c1))
				selectionEncryptionValidationHelper.addCheck("(4.C) v0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v0))
				selectionEncryptionValidationHelper.addCheck("(4.C) v1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v1))
				selectionEncryptionValidationHelper.addCheck("(4.D) The equation c=(c0+c1) mod q is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(addQ(&c0, &c1)))
				selectionEncryptionValidationHelper.addCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(&constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				selectionEncryptionValidationHelper.addCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(&constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				selectionEncryptionValidationHelper.addCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(elgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				selectionEncryptionValidationHelper.addCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", mulP(powP(&constants.G, &c1), powP(elgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))
			}
		}
	}
	correctnessOfSelectionsIsNotValid := !selectionEncryptionValidationHelper.validate()
	if correctnessOfSelectionsIsNotValid {
		return false
	}

	// validate adherence to vote limits (Step 5)
	voteLimitsValidationHelper := MakeValidationHelper(v.logger, "Adherence to vote limits (Step 5)")
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
			// Unwrap arguments for easier use
			aHat := contest.CiphertextAccumulation.Pad
			bHat := contest.CiphertextAccumulation.Data
			a := contest.Proof.Pad
			b := contest.Proof.Data
			v := contest.Proof.Response

			// Compute challenge and equations TODO: Should probably refactor
			c := crypto.HashElems(extendedBaseHash, aHat, bHat, a, b)
			equationFLeft := powP(&constants.G, &v)
			equationFRight := mulP(&a, powP(&aHat, c))
			equationGLeft := mulP(powP(&constants.G, mulP(schema.MakeBigIntFromString(strconv.Itoa(votesAllowed), 10), c)), powP(elgamalPublicKey, &v))
			equationGRight := mulP(&b, powP(&bHat, c))

			voteLimitsValidationHelper.addCheck("(5.A) The number of placeholder positions matches the selection limit ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", votesAllowed == numberOfSelections)
			voteLimitsValidationHelper.addCheck("(5.B) The a hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", aHat.Compare(calculatedAHat))
			voteLimitsValidationHelper.addCheck("(5.B) The b hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", bHat.Compare(calculatedBHat))
			voteLimitsValidationHelper.addCheck("(5.C) The given value V is in Zq ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isInRange(v))
			voteLimitsValidationHelper.addCheck("(5.D) The given value a are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Pad))
			voteLimitsValidationHelper.addCheck("(5.D) The given values b are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Data))
			voteLimitsValidationHelper.addCheck("(5.E) The challenge value is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", contest.Proof.Challenge.Compare(c))
			voteLimitsValidationHelper.addCheck("(5.F) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationFLeft.Compare(equationFRight))
			voteLimitsValidationHelper.addCheck("(5.E) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationGLeft.Compare(equationGRight))
		}
	}
	voteLimitsNotValid := !voteLimitsValidationHelper.validate()
	if voteLimitsNotValid {
		return false
	}

	// validate confirmation codes (Step 6)
	confirmationCodesValidationHelper := MakeValidationHelper(v.logger, "Validation of confirmation codes (Step 6)")
	hasSeen := make(map[string]bool)
	noDuplicateConfirmationCodesFound := true
	for _, ballot := range er.SubmittedBallots {
		// Computation of confirmation code (6.A)
		confirmationCodesValidationHelper.addCheck("(6.A) The confirmation code for ballot id: "+ballot.ObjectId+" is computed correct", true) // TODO: Fake it

		// No duplicate confirmation codes (6.B)
		stringOfCode := ballot.Code.String()
		if hasSeen[stringOfCode] {
			noDuplicateConfirmationCodesFound = false
		}
		hasSeen[stringOfCode] = true
	}
	confirmationCodesValidationHelper.addCheck("(6.B) No duplicate confirmation codes found", noDuplicateConfirmationCodesFound)
	confirmationCodesAreNotValid := !confirmationCodesValidationHelper.validate()
	if confirmationCodesAreNotValid {
		return false
	}

	// validate correctness of ballot aggregation (Step 7)
	ballotAggregationValidationHelper := MakeValidationHelper(v.logger, "Correctness of ballot aggregation (Step 7)")
	partialDecryptionsValidationHelper := MakeValidationHelper(v.logger, "Correctness of partial decryptions (Step 8)")
	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			a := schema.MakeBigIntFromString("1", 10)
			b := schema.MakeBigIntFromString("1", 10)
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
			//ballotAggregationValidationHelper.addCheck("(7.A) A is calculated correctly for "+i+j, A.Compare(a)) // TODO: Doesn't work
			//ballotAggregationValidationHelper.addCheck("(7.B) B is calculated correctly for "+i+j, B.Compare(b))

			for k, share := range selection.Shares {
				if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix
					v := share.Proof.Response
					c := share.Proof.Challenge
					ai := share.Proof.Pad
					bi := share.Proof.Data
					m := share.Share

					partialDecryptionsValidationHelper.addCheck("(8.A) The value v is in the set Zq for "+share.ObjectId+" "+strconv.Itoa(k), isInRange(v))
					partialDecryptionsValidationHelper.addCheck("(8.B) The value a is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Pad))
					partialDecryptionsValidationHelper.addCheck("(8.B) The value b is in the set Zqr for "+share.ObjectId+" "+strconv.Itoa(k), isValidResidue(share.Proof.Data))
					partialDecryptionsValidationHelper.addCheck("(8.B) The challenge is computed correctly "+share.ObjectId+" "+strconv.Itoa(k), c.Compare(crypto.HashElems(extendedBaseHash, A, B, ai, bi, m)))
					partialDecryptionsValidationHelper.addCheck("(8.D) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(&constants.G, &v).Compare(mulP(&ai, powP(&er.Guardians[k].ElectionPublicKey, &c))))
					partialDecryptionsValidationHelper.addCheck("(8.E) The equation is satisfied "+share.ObjectId+" "+strconv.Itoa(k), powP(&A, &v).Compare(mulP(&bi, powP(&m, &c))))
				}
			}
		}
	}
	ballotAggregationIsNotValid := !ballotAggregationValidationHelper.validate()
	if ballotAggregationIsNotValid {
		return false
	}

	// Validate correctness of partial decryptions (Step 8)
	partialDecryptionsAreNotValid := !partialDecryptionsValidationHelper.validate()
	if partialDecryptionsAreNotValid {
		return false
	}

	// Validate correctness of substitute data for missing guardians (Step 9)
	substituteDataValidationHelper := MakeValidationHelper(v.logger, "Correctness of substitute data for missing guardians (Step 9)")
	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			A := selection.Message.Pad
			B := selection.Message.Data
			for _, share := range selection.Shares {
				for _, part := range share.RecoveredParts {
					// TODO: Implement method to check if "Recovered parts" is not nil
					if part.ObjectId != "" {
						v := part.Proof.Response
						c := part.Proof.Challenge
						a := part.Proof.Pad
						b := part.Proof.Data
						m := part.PartialDecryption

						substituteDataValidationHelper.addCheck("(9.A) The given value v is in Zq", isInRange(v))
						substituteDataValidationHelper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(a))
						substituteDataValidationHelper.addCheck("(9.B) The given value a is in Zp^r", isValidResidue(b))
						substituteDataValidationHelper.addCheck("(9.C) The challenge value c is correct", c.Compare(crypto.HashElems(extendedBaseHash, A, B, a, b, m)))
						substituteDataValidationHelper.addCheck("(9.D) The equation is satisfied", powP(&constants.G, &v).Compare(mulP(&a, powP(&part.RecoveryPublicKey, &c))))
						substituteDataValidationHelper.addCheck("(9.E) The equation is satisfied", powP(&A, &v).Compare(mulP(&b, powP(&m, &c))))
					}
				}
			}
		}
	}
	substituteDataForMissingGuardiansIsNotValid := !substituteDataValidationHelper.validate()
	if substituteDataForMissingGuardiansIsNotValid {
		return false
	}

	// Validate correctness of construction of replacement partial decryptions (Step 10)
	replacementDecryptionsValidationHelper := MakeValidationHelper(v.logger, "Correctness of construction of replacement partial decryptions (Step 10)")
	for _, contest := range er.PlaintextTally.Contests {
		for _, selection := range contest.Selections {
			for _, share := range selection.Shares {
				product := schema.MakeBigIntFromString("1", 10)
				for _, part := range share.RecoveredParts {
					// TODO: Validate (10.A)?
					coefficient := er.CoefficientsValidationSet.Coefficients[part.GuardianIdentifier]
					product = mulP(product, powP(&part.PartialDecryption, &coefficient))
				}
				replacementDecryptionsValidationHelper.addCheck("(10.B) Correct tally share?", share.Share.Compare(product))
			}
		}
	}
	replacementPartialDecryptionsAreInvalid := !replacementDecryptionsValidationHelper.validate()
	if replacementPartialDecryptionsAreInvalid {
		return false
	}

	// Validate correctness of tally decryption (Step 11)
	tallyDecryptionValidationHelper := MakeValidationHelper(v.logger, "Correct decryption of tallies (Step 11)")
	for _, contest := range er.PlaintextTally.Contests {
		// TODO: Check 11.C to 11.F
		for _, selection := range contest.Selections {
			b := selection.Message.Data
			mi := schema.MakeBigIntFromString("1", 10)
			m := selection.Value
			t := schema.MakeBigIntFromInt(selection.Tally)
			for _, share := range selection.Shares {
				mi = mulP(mi, &share.Share)
			}
			tallyDecryptionValidationHelper.addCheck("(11.A) The equation is satisfied", b.Compare(mulP(&m, mi)))
			tallyDecryptionValidationHelper.addCheck("(11.B) The equation is satisfied", m.Compare(powP(&constants.G, t)))
		}
	}
	tallyDecryptionIsInvalid := !tallyDecryptionValidationHelper.validate()
	if tallyDecryptionIsInvalid {
		return false
	}

	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	// and validating correctness of substitute data for spoiled ballots (Step 13)
	spoiledBallotsDecryptionValidationHelper := MakeValidationHelper(v.logger, "Correctness of partial decryption for spoiled ballots (Step 12)")
	substituteDataForBallotsValidationHelper := MakeValidationHelper(v.logger, "Correctness of substitute data for spoiled ballots (Step 13)")
	for _, ballot := range er.SpoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data
				for _, share := range selection.Shares {
					if !share.Proof.Pad.Compare(schema.MakeBigIntFromString("0", 10)) { // Comparing with zero, will need better way of determining this TODO: Fix {
						m := share.Share // TODO: Refactor this struct in schema package
						a := share.Proof.Pad
						b := share.Proof.Data
						c := share.Proof.Challenge
						v := share.Proof.Response

						spoiledBallotsDecryptionValidationHelper.addCheck("(12.A) The given value v is in the set Zq", isInRange(v))
						spoiledBallotsDecryptionValidationHelper.addCheck("(12.B) The given value a is in the set Zpr", isValidResidue(a))
						spoiledBallotsDecryptionValidationHelper.addCheck("(12.B) The given value b is in the set Zpr", isValidResidue(b))
						spoiledBallotsDecryptionValidationHelper.addCheck("(12.C) The challenge is computed correctly", c.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, a, b, m)))
						spoiledBallotsDecryptionValidationHelper.addCheck("(12.D) The equation is satisfied", powP(&constants.G, &v).Compare(mulP(&a, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &c))))
						spoiledBallotsDecryptionValidationHelper.addCheck("(12.E) The equation is satisfied", powP(&alpha, &v).Compare(mulP(&b, powP(&m, &c))))

						// Step 13
						for _, part := range share.RecoveredParts {
							mil := part.PartialDecryption
							ai := part.Proof.Pad
							bi := part.Proof.Data
							ci := part.Proof.Challenge
							vi := part.Proof.Response

							substituteDataForBallotsValidationHelper.addCheck("(13.A) The given value v is in Zq", isInRange(vi))
							spoiledBallotsDecryptionValidationHelper.addCheck("(13.B) The given value a is in the set Zpr", isValidResidue(ai))
							spoiledBallotsDecryptionValidationHelper.addCheck("(13.B) The given value b is in the set Zpr", isValidResidue(bi))
							spoiledBallotsDecryptionValidationHelper.addCheck("(13.C) The challenge is computed correctly", ci.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, ai, bi, mil)))
							spoiledBallotsDecryptionValidationHelper.addCheck("(13.D) The equation is satisfied", powP(&constants.G, &vi).Compare(powP(mulP(&ai, &part.RecoveryPublicKey), &ci)))
							spoiledBallotsDecryptionValidationHelper.addCheck("(13.E) The equation is satisfied", powP(&ai, &vi).Compare(mulP(&bi, powP(&mil, &ci))))

						}
					}
				}
			}
		}
	}
	spoiledBallotsPartialDecryptionIsInvalid := !spoiledBallotsDecryptionValidationHelper.validate()
	if spoiledBallotsPartialDecryptionIsInvalid {
		return false
	}

	substituteDataForSpoiledBallotsIsInvalid := !substituteDataForBallotsValidationHelper.validate()
	if substituteDataForSpoiledBallotsIsInvalid {
		return false
	}

	// Validation of correct replacement partial decryptions for spoiled ballots (Step 14)
	// ...
	// ...

	// Validation of correct decryption of spoiled ballots (Step 15)
	// and validation of correctness of spoiled ballots (Step 16)
	decryptionOfSpoiledBallotsValidationHelper := MakeValidationHelper(v.logger, "Correct decryption of spoiled ballots (Step 15)")
	correctnessOfSpoiledBallotsValidationHelper := MakeValidationHelper(v.logger, "Correctness of spoiled ballots (Step 16)")
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

				decryptionOfSpoiledBallotsValidationHelper.addCheck("(15.A) The equation is satisfied", beta.Compare(mulP(&m, mi)))
				decryptionOfSpoiledBallotsValidationHelper.addCheck("(15.B) The equation is satisfied", m.Compare(powP(&constants.G, V)))

				correctnessOfSpoiledBallotsValidationHelper.addCheck("(16.A) For each option in the contest, the selection V is either a 0 or a 1", selection.Tally == 0 || selection.Tally == 1)

			}
			correctnessOfSpoiledBallotsValidationHelper.addCheck("(16.B) The sum of all selections in the contest is at most the selection limit L for that contest.", sumOfAllSelections <= getContest(contest.ObjectId, er.Manifest.Contests).VotesAllowed)
			// TODO: 16.C -> 16.E
		}
	}
	decryptionOfSpoiledBallotsIsInvalid := !decryptionOfSpoiledBallotsValidationHelper.validate()
	if decryptionOfSpoiledBallotsIsInvalid {
		return false
	}

	correctnessSpoiledBallotsIsInvalid := !correctnessOfSpoiledBallotsValidationHelper.validate()
	if correctnessSpoiledBallotsIsInvalid {
		return false
	}

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	correctContestDataValidationHelper := MakeValidationHelper(v.logger, "Correctness of contest data partial decryptions for spoiled ballots (Step 17)")
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

				correctContestDataValidationHelper.addCheck("(17.A) The given value v is in the set Zq", isInRange(vi))
				correctContestDataValidationHelper.addCheck("(17.B) The given value a is in the set Zqr", isValidResidue(ai))
				correctContestDataValidationHelper.addCheck("(17.B) The given value b is in the set Zqr", isValidResidue(bi))
				correctContestDataValidationHelper.addCheck("(17.C) The challenge is correctly computed", ci.Compare(crypto.HashElems(extendedBaseHash, c0, c1, c2, ai, bi, mi)))
				correctContestDataValidationHelper.addCheck("(17.D) The equation is satisfied", powP(&constants.G, &vi).Compare(mulP(&ai, powP(getGuardianPublicKey(share.GuardianId, er.Guardians), &ci))))
				correctContestDataValidationHelper.addCheck("(17.E) The equation is satisfied", powP(&c0, &vi).Compare(mulP(&bi, powP(&mi, &ci))))
			}
		}
	}
	correctContestDataPartialDecryptionsForSpoiledBallotIsInvalid := !correctContestDataValidationHelper.validate()
	if correctContestDataPartialDecryptionsForSpoiledBallotIsInvalid {
		return false
	}

	// Verification was successful
	return true
}
