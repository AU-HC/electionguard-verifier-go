package core

import (
	"electionguard-verifier-go/deserialize"
)

func MakeVerifyStrategy(useMultipleThreads bool, amountOfLogicalCores int) VerifyStrategy {
	if useMultipleThreads {
		return ParallelStrategy{amountOfLogicalCores: amountOfLogicalCores}
	}

	return SingleThreadStrategy{}
}

type VerifyStrategy interface {
	verify(er *deserialize.ElectionRecord, verifier *Verifier)
	getBallotChunkSize(amountOfBallots int) int
	getContestSplitSize() int
}

type SingleThreadStrategy struct {
}

func (s SingleThreadStrategy) verify(er *deserialize.ElectionRecord, verifier *Verifier) {
	// Validate election parameters (Step 1)
	verifier.validateElectionConstants(er)

	// Validate guardian public-key (Step 2)
	verifier.validateGuardianPublicKeys(er)

	// Validation election public-key (Step 3)
	verifier.validateJointPublicKey(er)

	// Validate correctness of selection encryptions (Step 4)
	verifier.validateSelectionEncryptions(er)

	// Validate adherence to vote limits (Step 5)
	verifier.validateVoteLimits(er)

	// Validate confirmation codes (Step 6)
	verifier.validateConfirmationCodes(er)

	// Validate correctness of ballot aggregation (Step 7)
	verifier.validateBallotAggregation(er)

	// Validate correctness of partial decryptions (Step 8)
	verifier.validatePartialDecryptions(er)

	// Validate correctness of substitute data for missing guardians (Step 9)
	verifier.validateSubstituteDataForMissingGuardians(er)

	// Validate correctness of construction of replacement partial decryptions (Step 10)
	verifier.validateConstructionOfReplacementForPartialDecryptions(er)

	// Validate correctness of tally decryption (Step 11)
	verifier.validateTallyDecryption(er)

	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	verifier.validatePartialDecryptionForSpoiledBallots(er)

	// Validate correctness of substitute data for spoiled ballots (Step 13)
	verifier.validateSubstituteDataForSpoiledBallots(er)

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)
	verifier.validateReplacementPartialDecryptionForSpoiledBallots(er)

	// Validation of correct decryption of spoiled ballots (Step 15)
	verifier.validateDecryptionOfSpoiledBallots(er)

	// and validation of correctness of spoiled ballots (Step 16)
	verifier.validateCorrectnessOfSpoiledBallots(er)

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	verifier.validateContestDataPartialDecryptionsForSpoiledBallots(er)

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	verifier.validateSubstituteContestDataForSpoiledBallots(er)

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)
	verifier.validateContestReplacementDecryptionForSpoiledBallots(er)
}

func (s SingleThreadStrategy) getBallotChunkSize(amountOfBallots int) int {
	return amountOfBallots
}

func (s SingleThreadStrategy) getContestSplitSize() int {
	return 1
}

type ParallelStrategy struct {
	amountOfLogicalCores int // Amount of logical cores on the current machine, used to decide amount of goroutines
}

func (s ParallelStrategy) verify(er *deserialize.ElectionRecord, verifier *Verifier) {
	// Validate election parameters (Step 1)
	go verifier.validateElectionConstants(er)

	// Validate guardian public-key (Step 2)
	go verifier.validateGuardianPublicKeys(er)

	// Validation election public-key (Step 3)
	go verifier.validateJointPublicKey(er)

	// Validate correctness of selection encryptions (Step 4)
	go verifier.validateSelectionEncryptions(er)

	// Validate adherence to vote limits (Step 5)
	go verifier.validateVoteLimits(er)

	// Validate confirmation codes (Step 6)
	go verifier.validateConfirmationCodes(er)

	// Validate correctness of ballot aggregation (Step 7)
	go verifier.validateBallotAggregation(er)

	// Validate correctness of partial decryptions (Step 8)
	go verifier.validatePartialDecryptions(er)

	// Validate correctness of substitute data for missing guardians (Step 9)
	go verifier.validateSubstituteDataForMissingGuardians(er)

	// Validate correctness of construction of replacement partial decryptions (Step 10)
	go verifier.validateConstructionOfReplacementForPartialDecryptions(er)

	// Validate correctness of tally decryption (Step 11)
	go verifier.validateTallyDecryption(er)

	// Validate correctness of partial decryption for spoiled ballots (Step 12)
	go verifier.validatePartialDecryptionForSpoiledBallots(er)

	// Validate correctness of substitute data for spoiled ballots (Step 13)
	go verifier.validateSubstituteDataForSpoiledBallots(er)

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)
	go verifier.validateReplacementPartialDecryptionForSpoiledBallots(er)

	// Validation of correct decryption of spoiled ballots (Step 15)
	go verifier.validateDecryptionOfSpoiledBallots(er)

	// and validation of correctness of spoiled ballots (Step 16)
	go verifier.validateCorrectnessOfSpoiledBallots(er)

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)
	go verifier.validateContestDataPartialDecryptionsForSpoiledBallots(er)

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)
	go verifier.validateSubstituteContestDataForSpoiledBallots(er)

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)
	go verifier.validateContestReplacementDecryptionForSpoiledBallots(er)

	// Waiting for all goroutines to finish
	verifier.wg.Wait()
}

func (s ParallelStrategy) getBallotChunkSize(amountOfBallots int) int {
	if amountOfBallots > s.amountOfLogicalCores {
		return amountOfBallots / s.amountOfLogicalCores
	}
	return amountOfBallots
}

func (s ParallelStrategy) getContestSplitSize() int {
	return 2
}
