package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
	"go.uber.org/zap"
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
	er, electionRecordIsValid := v.getElectionRecord(path)
	if !electionRecordIsValid {
		return false
	}

	// Validate election parameters (Step 1)
	v.constants = utility.MakeCorrectElectionConstants()
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

func (v *Verifier) getElectionRecord(path string) (*deserialize.ElectionRecord, bool) {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	er, err := parser.ParseElectionRecord(path)

	if err != "" {
		v.logger.Info("[INVALID]: Election data was well formed (Step 0)")
		v.logger.Debug(err)
	} else {
		v.logger.Info("[VALID]: Election data was well formed (Step 0)")
	}

	// If length of error message is 0, no errors were reported
	return er, len(err) == 0
}
