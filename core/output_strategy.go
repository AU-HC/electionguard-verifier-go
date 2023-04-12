package core

import (
	"electionguard-verifier-go/deserialize"
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
	ElectionName      string
	VerificationSteps []ValidationHelper
}

type OutputStrategy interface {
	Output(record deserialize.ElectionRecord, results []*ValidationHelper)
	OutputBenchmark(amountOfSamples int, runs []float64)
}

type NoOutputStrategy struct {
}

func (s NoOutputStrategy) Output(record deserialize.ElectionRecord, results []*ValidationHelper) {
	// do nothing
}

func (s NoOutputStrategy) OutputBenchmark(samples int, runs []float64) {
	// do nothing
}

type ToFileStrategy struct {
	Path string
}

func (s ToFileStrategy) Output(record deserialize.ElectionRecord, results []*ValidationHelper) {
	var xd []ValidationHelper
	for i, xd2 := range results {
		if i != 0 {
			xd = append(xd, *xd2)
		}
	}
	vr := VerificationRecord{ElectionName: record.Manifest.ElectionScopeID, VerificationSteps: xd}
	jsonBytes, err := json.MarshalIndent(vr, "", "  ")
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
