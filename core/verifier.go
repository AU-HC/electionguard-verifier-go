package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/utility"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Verifier struct {
	logger           *zap.Logger                      // logger used to log information
	constants        utility.CorrectElectionConstants // constants is election constants (p, q, r, g)
	wg               *sync.WaitGroup                  // wg is used to sync goroutines for each step
	helpers          []*ValidationHelper              // helpers are used to store result of each verification step
	verifierStrategy VerifyStrategy                   // verifierStrategy is used to decide if the steps should be verified concurrently
	outputStrategy   OutputStrategy                   // outputStrategy is used to output the verification results
}

func MakeVerifier(logger *zap.Logger) *Verifier {
	return &Verifier{logger: logger, wg: &sync.WaitGroup{}, helpers: make([]*ValidationHelper, 20)}
}

func (v *Verifier) Verify(path string) bool {
	// Deserialize election record and fetching correct election constants (Step 0)
	er, electionRecordIsNotValid := v.getElectionRecord(path)
	v.constants = utility.MakeCorrectElectionConstants()
	if electionRecordIsNotValid {
		return false
	}

	// Setting up synchronization (Will have to even if using one thread)
	v.wg.Add(19)

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
	return electionRecord, len(err) != 0
}

func (v *Verifier) SetOutputStrategy(strategy OutputStrategy) {
	v.outputStrategy = strategy
}

func (v *Verifier) SetVerifyStrategy(strategy VerifyStrategy) {
	v.verifierStrategy = strategy
}
