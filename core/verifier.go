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
	args := parser.ConvertJsonDataToGoStruct(path)
	v.logger.Info("[VALID]: Election data was formed well (Step 0)")

	// Validate election parameters (Step 1):
	constants := utility.MakeCorrectElectionConstants()
	electionParametersHelper := MakeValidationHelper(v.logger, "Election parameters (Step 1)")
	electionParametersHelper.AddCheck("(1.A) The large prime is equal to the large modulus p", constants.P.Compare(&args.ElectionConstants.LargePrime))
	electionParametersHelper.AddCheck("(1.B) The small prime is equal to the prime q", constants.Q.Compare(&args.ElectionConstants.SmallPrime))
	electionParametersHelper.AddCheck("(1.C) The cofactor is equal to r = (p − 1)/q", constants.C.Compare(&args.ElectionConstants.Cofactor))
	electionParametersHelper.AddCheck("(1.D) The generator is equal to the generator g", constants.G.Compare(&args.ElectionConstants.Generator))
	electionParametersIsNotValid := !electionParametersHelper.Validate()
	if electionParametersIsNotValid {
		return false
	}

	// Validate guardian public-key (Step 2)
	publicKeyValidationHelper := MakeValidationHelper(v.logger, "Guardian public-key validation (Step 2)")
	electionKeyValidationHelper := MakeValidationHelper(v.logger, "Election public-key validation (Step 3)")
	elgamalPublicKey := schema.MakeBigIntFromString("1", 10)
	for i, guardian := range args.Guardians {
		elgamalPublicKey = mulP(elgamalPublicKey, &guardian.ElectionPublicKey)
		for j, proof := range guardian.ElectionProofs {
			// (2.A)
			hash := crypto.HashElems(guardian.ElectionCommitments[j], proof.Commitment)
			publicKeyValidationHelper.AddCheck("(2.A) The challenge is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", proof.Challenge.Compare(hash))

			// (2.B)
			left := powP(&constants.G, &proof.Response)
			right := mulP(powP(&guardian.ElectionCommitments[j], &proof.Challenge), &proof.Commitment)
			publicKeyValidationHelper.AddCheck("(2.B) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", left.Compare(right))
		}
	}
	publicKeysAreNotValid := !publicKeyValidationHelper.Validate()
	if publicKeysAreNotValid {
		return false
	}

	// Validate election public key (Step 3) [ERROR IN SPEC SHEET FOR (3.B)]
	extendedBaseHash := args.CiphertextElectionRecord.CryptoExtendedBaseHash
	computedExtendedBaseHash := crypto.HashElems(args.CiphertextElectionRecord.CryptoBaseHash, args.CiphertextElectionRecord.CommitmentHash)

	electionKeyValidationHelper.AddCheck("(3.A) The joint public election key is computed correctly", elgamalPublicKey.Compare(&args.CiphertextElectionRecord.ElgamalPublicKey))
	electionKeyValidationHelper.AddCheck("(3.B) The extended base hash is computed correctly", extendedBaseHash.Compare(computedExtendedBaseHash))

	jointElectionKeyIsNotValid := !electionKeyValidationHelper.Validate()
	if jointElectionKeyIsNotValid {
		return false
	}

	// Validate correctness of selection encryptions (Step 4)
	selectionEncryptionValidationHelper := MakeValidationHelper(v.logger, "Correctness of selection encryptions (Step 4)")
	for i, ballot := range args.SubmittedBallots {
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
				selectionEncryptionValidationHelper.AddCheck("(4.A) a is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a))
				selectionEncryptionValidationHelper.AddCheck("(4.A) b is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b))
				selectionEncryptionValidationHelper.AddCheck("(4.A) a0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a0))
				selectionEncryptionValidationHelper.AddCheck("(4.A) b0 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b0))
				selectionEncryptionValidationHelper.AddCheck("(4.A) a1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(a1))
				selectionEncryptionValidationHelper.AddCheck("(4.A) b1 is in the set Z_p^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isValidResidue(b1))
				selectionEncryptionValidationHelper.AddCheck("(4.B) The challenge value c is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(crypto.HashElems(args.CiphertextElectionRecord.CryptoExtendedBaseHash, a, b, a0, b0, a1, b1)))
				selectionEncryptionValidationHelper.AddCheck("(4.C) c0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c0))
				selectionEncryptionValidationHelper.AddCheck("(4.C) c1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(c1))
				selectionEncryptionValidationHelper.AddCheck("(4.C) v0 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v0))
				selectionEncryptionValidationHelper.AddCheck("(4.C) v1 is in Zq for ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", isInRange(v1))
				selectionEncryptionValidationHelper.AddCheck("(4.D) The equation c=(c0+c1) mod q is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", c.Compare(addQ(&c0, &c1)))
				selectionEncryptionValidationHelper.AddCheck("(4.E) The equation g^v0=a0*a^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(&constants.G, &v0).Compare(mulP(&a0, powP(&a, &c0))))
				selectionEncryptionValidationHelper.AddCheck("(4.F) The equation g^v1=a1*a^c1 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(&constants.G, &v1).Compare(mulP(&a1, powP(&a, &c1))))
				selectionEncryptionValidationHelper.AddCheck("(4.G) The equation K^v0=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", powP(elgamalPublicKey, &v0).Compare(mulP(&b0, powP(&b, &c0))))
				selectionEncryptionValidationHelper.AddCheck("(4.H) The equation g^c1=b0*b^c0 is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+","+strconv.Itoa(k)+")", mulP(powP(&constants.G, &c1), powP(elgamalPublicKey, &v1)).Compare(mulP(&b1, powP(&b, &c1))))
			}
		}
	}
	correctnessOfSelectionsIsNotValid := !selectionEncryptionValidationHelper.Validate()
	if correctnessOfSelectionsIsNotValid {
		return false
	}

	// Validate adherence to vote limits (Step 5)
	voteLimitsValidationHelper := MakeValidationHelper(v.logger, "Adherence to vote limits (Step 5)")
	for i, ballot := range args.SubmittedBallots {
		for j, contest := range ballot.Contests {
			contestInManifest := findContestFromObjectID(contest.ObjectId, args.Manifest.Contests)
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

			voteLimitsValidationHelper.AddCheck("(5.A) The number of placeholder positions matches the selection limit ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", votesAllowed == numberOfSelections)
			voteLimitsValidationHelper.AddCheck("(5.B) The a hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", aHat.Compare(calculatedAHat))
			voteLimitsValidationHelper.AddCheck("(5.B) The b hat is computed correctly ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", bHat.Compare(calculatedBHat))
			voteLimitsValidationHelper.AddCheck("(5.C) The given value V is in Zq ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isInRange(v))
			voteLimitsValidationHelper.AddCheck("(5.D) The given value a are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Pad))
			voteLimitsValidationHelper.AddCheck("(5.D) The given values b are in Zp^r ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", isValidResidue(contest.Proof.Data))
			voteLimitsValidationHelper.AddCheck("(5.E) The challenge value is correctly computed ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", contest.Proof.Challenge.Compare(c))
			voteLimitsValidationHelper.AddCheck("(5.F) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationFLeft.Compare(equationFRight))
			voteLimitsValidationHelper.AddCheck("(5.E) The equation is satisfied ("+strconv.Itoa(i)+","+strconv.Itoa(j)+")", equationGLeft.Compare(equationGRight))
		}
	}
	voteLimitsNotValid := !voteLimitsValidationHelper.Validate()
	if voteLimitsNotValid {
		return false
	}

	// Validate confirmation codes (Step 6)
	// ...
	// ...

	// Verification was successful
	return true
}

func findContestFromObjectID(objectID string, contests []schema.Contest) schema.Contest {
	for _, contest := range contests {
		if objectID == contest.ObjectID {
			return contest
		}
	}

	return schema.Contest{}
}
