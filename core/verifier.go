package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

var amountOfVerificationSteps = 13

type Verifier struct {
	logger           *zap.Logger              // logger used to log information
	constants        CorrectElectionConstants // constants is election constants (p, q, r, g)
	wg               *sync.WaitGroup          // wg is used to sync goroutines for each step
	helpers          []*ValidationHelper      // helpers are used to store result of each verification step
	verifierStrategy VerifyStrategy           // verifierStrategy is used to decide if the steps should be verified concurrently
	outputStrategy   OutputStrategy           // outputStrategy is used to output the verification results
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{
		logger:    logger,
		wg:        &sync.WaitGroup{},
		helpers:   make([]*ValidationHelper, amountOfVerificationSteps+1),
		constants: MakeCorrectElectionConstants(),
	}
}

func (v *Verifier) Verify(path string) bool {
	// Deserialize election record and fetching correct election constants (Step 0)
	er, electionRecordIsNotValid := v.getElectionRecord(path)
	if electionRecordIsNotValid {
		return false
	}

	// Setting up synchronization (will have to even if using one thread)
	v.wg.Add(amountOfVerificationSteps)

	// Starting time and verifying election using supplied strategy
	start := time.Now()
	v.verifierStrategy.verify(er, v)
	elapsed := time.Since(start)

	electionIsValid := v.validateAllVerificationSteps()
	v.logger.Info("Validation of election took: " + elapsed.String())

	// Output validation results to file using specific strategy
	v.outputStrategy.Output(*er, v.helpers)

	return electionIsValid
}

func (v *Verifier) validateAllVerificationSteps() bool {
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

func (v *Verifier) getElectionRecord(path string) (*schema.ElectionRecord, bool) {
	// Fetch and deserialize election data (Step 0)
	parser := *deserialize.MakeParser(v.logger)
	electionRecord, err := parser.ParseElectionRecord(path)

	if err == "" {
		v.logger.Info("Election data was well formed")
	} else {
		v.logger.Info("Election data was not well formed")
		v.logger.Debug(err)
	}

	// If length of error message is 0, no errors were reported and thus return electionRecord, true
	return electionRecord, len(err) != 0
}

func (v *Verifier) SetOutputStrategy(strategy OutputStrategy) {
	v.outputStrategy = strategy
}

func (v *Verifier) SetVerifyStrategy(strategy VerifyStrategy) {
	v.verifierStrategy = strategy
}

func (v *Verifier) Benchmark(path string, amountOfSamples int) {
	er, electionRecordIsNotValid := v.getElectionRecord(path)
	if electionRecordIsNotValid {
		return
	}

	runs := make([]float64, amountOfSamples)
	for i := 0; i < amountOfSamples; i++ {
		// Setting up synchronization (Will have to even if using one thread)
		v.wg.Add(13)

		// Starting time and verifying election using supplied strategy
		start := time.Now()
		v.verifierStrategy.verify(er, v)
		elapsed := time.Since(start)

		runs[i] = float64(elapsed.Milliseconds())

		v.logger.Info(strconv.Itoa(i+1) + "/" + strconv.Itoa(amountOfSamples))
	}

	// Output the data to a json file
	v.outputStrategy.OutputBenchmark(amountOfSamples, runs)
}

func (v *Verifier) BenchmarkDeserialization(amountOfSamples int) {
	paths := make([]string, 8)
	paths[0] = "data/sandbox_25/election_record/"
	paths[1] = "data/sandbox_50/election_record/"
	paths[2] = "data/sandbox_100/election_record/"
	paths[3] = "data/sandbox_250/election_record/"
	paths[4] = "data/sandbox_500/election_record/"
	paths[4] = "data/sandbox_750/election_record/"
	paths[5] = "data/sandbox_1000/election_record/"
	paths[6] = "data/sandbox/election_record/"

	for _, path := range paths {
		var totalTime int64
		for i := 0; i < amountOfSamples; i++ {
			start := time.Now()
			_, _ = v.getElectionRecord(path)
			totalTime += time.Since(start).Milliseconds()
		}

		totalTimeInMs := int(totalTime) / amountOfSamples
		fmt.Println(fmt.Sprintf("\"%s\" had mean time of: %d ms", path, totalTimeInMs))
	}
}
