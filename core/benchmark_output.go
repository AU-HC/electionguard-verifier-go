package core

import (
	"encoding/json"
	"github.com/montanaflynn/stats"
	"os"
)

type BenchmarkResult struct {
	AmountOfSamples int
	MinRun          float64
	MaxRun          float64
	Median          float64
	Mean            float64
	Runs            []float64
}

func MakeBenchmarkResults(samples int, runs []float64) BenchmarkResult {
	// Create struct
	benchmarkResults := BenchmarkResult{AmountOfSamples: samples, Runs: runs}

	// Fill the struct with data
	min, _ := stats.Min(runs)
	benchmarkResults.MinRun = min

	max, _ := stats.Max(runs)
	benchmarkResults.MaxRun = max

	median, _ := stats.Median(runs)
	benchmarkResults.Median = median

	mean, _ := stats.Mean(runs)
	benchmarkResults.Mean = mean

	return benchmarkResults
}

func (b *BenchmarkResult) OutputToJsonFile() {
	jsonBytes, err := json.MarshalIndent(*b, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("benchmark.json", jsonBytes, 0644)
	if err != nil {
		return
	}
}
