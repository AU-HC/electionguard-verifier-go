package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
	"go.uber.org/zap"
	"sync"
)

type Verifier struct {
	logger    *zap.Logger                      // logger used to log information
	constants utility.CorrectElectionConstants // constants is election constants (p, q, r, g)
	wg        *sync.WaitGroup                  // wg is used to sync goroutines for each step
	helpers   []*ValidationHelper              // helpers are used to store result of each verification step
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger, wg: &sync.WaitGroup{}, helpers: make([]*ValidationHelper, 20)}
}

func (v *Verifier) Verify(path string) bool {
	// Deserialize election record and fetching correct election constants (Step 0)
	er, electionRecordIsValid := v.getElectionRecord(path)
	v.constants = utility.MakeCorrectElectionConstants()
	if !electionRecordIsValid {
		return false
	}

	// Setting up synchronization
	v.wg.Add(19)

	// Validate election parameters (Step 1)
	go v.validateElectionConstants(er)

	// Validate guardian public-key (Step 2)
	go v.validateGuardianPublicKeys(er)

	// Validation election public-key (Step 3)
	go v.validateJointPublicKey(er)

	// Validate correctness of selection encryptions (Step 4)
	go v.validateSelectionEncryptions(er)

	// Validate adherence to vote limits (Step 5)
	go v.validateVoteLimits(er)

	// Validate confirmation codes (Step 6)
	go v.validateConfirmationCodes(er)

	// Validate correctness of ballot aggregation (Step 7)
	go v.validateBallotAggregation(er)

	// Validate correctness of partial decryptions (Step 8)
	go v.validatePartialDecryptions(er)

	// Validate correctness of substitute data for missing guardians (Step 9)
	go v.validateSubstituteDataForMissingGuardians(er)

	// Validate correctness of construction of replacement partial decryptions (Step 10)
	go v.validateConstructionOfReplacementForPartialDecryptions(er)

	// Validate correctness of tally decryption (Step 11)
	go v.validateTallyDecryption(er)

	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	go v.validatePartialDecryptionForSpoiledBallots(er)

	// Validate correctness of substitute data for spoiled ballots (Step 13)
	go v.validateSubstituteDataForSpoiledBallots(er)

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)
	go v.validateReplacementPartialDecryptionForSpoiledBallots(er)

	// Validation of correct decryption of spoiled ballots (Step 15)
	go v.validateDecryptionOfSpoiledBallots(er)

	// and validation of correctness of spoiled ballots (Step 16)
	go v.validateCorrectnessOfSpoiledBallots(er)

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	go v.validateContestDataPartialDecryptionsForSpoiledBallots(er)

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	go v.validateSubstituteContestDataForSpoiledBallots(er)

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)
	go v.validateContestReplacementDecryptionForSpoiledBallots(er)

	// Waiting for all goroutines to finish
	electionIsValid := v.validateAllVerificationSteps()

	// Output validation results to file using specific strategy

	return electionIsValid
}

func (v *Verifier) validateAllVerificationSteps() bool {
	// Waiting for all goroutines to finish
	v.wg.Wait()

	// Checking each step
	electionIsValid := true
	for i, result := range v.helpers {
		if i != 0 {
			verificationStepIsNotValid := !result.validate()
			if verificationStepIsNotValid {
				electionIsValid = false
			}
		}
	}

	return electionIsValid
}

func (v *Verifier) getElectionRecord(path string) (*deserialize.ElectionRecord, bool) {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	electionRecord, err := parser.ParseElectionRecord(path)

	if err != "" {
		v.logger.Info("[INVALID]: Election data was well formed (Step 0)")
		v.logger.Debug(err)
	} else {
		v.logger.Info("[VALID]: Election data was well formed (Step 0)")
	}

	// If length of error message is 0, no errors were reported and thus return electionRecord, true
	return electionRecord, len(err) == 0
}
