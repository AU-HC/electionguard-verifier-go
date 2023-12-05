package core

import (
	"electionguard-verifier-go/schema"
	"encoding/json"
	"os"
)

func MakeOutputStrategy(outputPath string) OutputStrategy {
	if outputPath == "" {
		return NoOutputStrategy{}
	}

	return ToFileStrategy{Path: outputPath}
}

type VerificationRecord struct {
	ElectionName       string
	VerificationStatus bool
	VerificationSteps  []ValidationHelper
}

type OutputStrategy interface {
	Output(record schema.ElectionRecord, results []*ValidationHelper)
	OutputBenchmark(amountOfSamples int, runs []float64)
}

type NoOutputStrategy struct {
}

func (s NoOutputStrategy) Output(record schema.ElectionRecord, results []*ValidationHelper) {
	// do nothing
}

func (s NoOutputStrategy) OutputBenchmark(samples int, runs []float64) {
	// do nothing
}

type ToFileStrategy struct {
	Path string
}

func (s ToFileStrategy) Output(record schema.ElectionRecord, results []*ValidationHelper) {
	var helpers []ValidationHelper
	electionIsValid := true

	for i, helper := range results {
		if i != 0 {
			if !helper.isValid {
				electionIsValid = false
			}
			helpers = append(helpers, *helper)
		}
	}

	verificationRecord := VerificationRecord{ElectionName: record.Manifest.ElectionScopeID, VerificationStatus: electionIsValid, VerificationSteps: helpers}
	jsonBytes, err := json.MarshalIndent(verificationRecord, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(s.Path, jsonBytes, 0644)
	if err != nil {
		return
	}
}

func (s ToFileStrategy) OutputBenchmark(amountOfSamples int, runs []float64) {
	b := MakeBenchmarkResults(amountOfSamples, runs)

	jsonBytes, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(s.Path, jsonBytes, 0644)
	if err != nil {
		return
	}
}
