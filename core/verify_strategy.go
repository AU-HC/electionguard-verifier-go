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

	// Validation election public-key (Step 3)

	// Validate correctness of selection encryptions (Step 4)

	// Validate adherence to vote limits (Step 5)

	// Validate confirmation codes (Step 6)

	// Validate correctness of ballot aggregation (Step 7)

	// Validate correctness of partial decryptions (Step 8)

	// Validate correctness of substitute data for missing guardians (Step 9)

	// Validate correctness of construction of replacement partial decryptions (Step 10)

	// Validate correctness of tally decryption (Step 11)

	// Validate correctness of partial decryption for spoiled ballots (Step 12)

	// Validate correctness of substitute data for spoiled ballots (Step 13)

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)

	// Validation of correct decryption of spoiled ballots (Step 15)

	// and validation of correctness of spoiled ballots (Step 16)

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)
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

	// Validation election public-key (Step 3)

	// Validate correctness of selection encryptions (Step 4)

	// Validate adherence to vote limits (Step 5)

	// Validate confirmation codes (Step 6)

	// Validate correctness of ballot aggregation (Step 7)

	// Validate correctness of partial decryptions (Step 8)

	// Validate correctness of substitute data for missing guardians (Step 9)

	// Validate correctness of construction of replacement partial decryptions (Step 10)

	// Validate correctness of tally decryption (Step 11)

	// Validate correctness of partial decryption for spoiled ballots (Step 12)

	// Validate correctness of substitute data for spoiled ballots (Step 13)

	// Validate of correct replacement partial decryptions for spoiled ballots (Step 14)

	// Validation of correct decryption of spoiled ballots (Step 15)

	// and validation of correctness of spoiled ballots (Step 16)

	// Verifying correctness of contest data partial decryptions for spoiled ballots (Step 17)

	// Validating correctness of substitute contest data for spoiled ballots (Step 18)

	// Validating the correctness of contest replacement decryptions for spoiled ballots (Step 19)

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
