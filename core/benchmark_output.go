package core

import (
	"github.com/montanaflynn/stats"
)

type BenchmarkResult struct {
	AmountOfSamples   int
	MinRun            float64
	MaxRun            float64
	Median            float64
	Mean              float64
	StandardDeviation float64
	Runs              []float64
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

	sd, _ := stats.StandardDeviation(runs)
	benchmarkResults.StandardDeviation = sd

	return benchmarkResults
}
