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
	verifier.validateParameters(er)

	// Validate guardian public-key (Step 2)
	verifier.validateGuardianPublicKeys(er)

	// Validation election public-key (Step 3)
	verifier.validateElectionPublicKey(er)

	// Validate correctness base hash (Step 4)
	verifier.validateExtendedBaseHash(er)

	// Validate correctness of selection encryptions (Step 5)
	// verifier.validateSelectionEncryptions(er)

	// Validate adherence to vote limits (Step 6)
	verifier.validateAdherenceToVoteLimits(er)

	// Validate confirmation codes (Step 7)
	verifier.validateConfirmationCodes(er)

	// Validate correctness of ballot aggregation (Step 8)
	verifier.validateBallotAggregation(er)

	// Validate correctness of tally decryptions (Step 9)
	verifier.validateTallyDecryptions(er)

	// Validate correct decryption of tallies (Step 10)
	verifier.validateCorrectnessOfTallyDecryptions(er)

	// Correctness of decryptions of contest data (Step 11)

	// Correctness of decryptions for challenged ballots (Step 12)

	// Validation of correct decryption of challenged ballots (Step 13)

	// Verifications 14-18 should not be implemented
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
	go verifier.validateParameters(er)

	// Validate guardian public-key (Step 2)
	go verifier.validateGuardianPublicKeys(er)

	// Validation election public-key (Step 3)
	go verifier.validateElectionPublicKey(er)

	// Validate correctness base hash (Step 4)
	go verifier.validateExtendedBaseHash(er)

	// Validate correctness of selection encryptions (Step 5)
	go verifier.validateSelectionEncryptions(er)

	// Validate adherence to vote limits (Step 6)
	go verifier.validateAdherenceToVoteLimits(er)

	// Validate confirmation codes (Step 7)
	go verifier.validateConfirmationCodes(er)

	// Validate correctness of ballot aggregation (Step 8)
	go verifier.validateBallotAggregation(er)

	// Validate correctness of tally decryptions (Step 9)
	go verifier.validateTallyDecryptions(er)

	// Validate correct decryption of tallies (Step 10)
	go verifier.validateCorrectnessOfTallyDecryptions(er)

	// Correctness of decryptions of contest data (Step 11)

	// Correctness of decryptions for challenged ballots (Step 12)

	// Validation of correct decryption of challenged ballots (Step 13)

	// Verifications 14-18 should not be implemented

	// Waiting for all goroutines to finish
	verifier.wg.Wait()
}

func (s ParallelStrategy) getBallotChunkSize(amountOfBallots int) int {
	if amountOfBallots > s.amountOfLogicalCores {
		return amountOfBallots / s.amountOfLogicalCores
	}
	if amountOfBallots > 5 {
		// in order to ensure multithreading in smaller elections where more logical cores exist than ballots.
		return amountOfBallots / 3
	}
	return amountOfBallots
}

func (s ParallelStrategy) getContestSplitSize() int {
	return 2
}
