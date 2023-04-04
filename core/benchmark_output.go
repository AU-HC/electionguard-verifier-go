package core

type BenchmarkResult struct {
	AmountOfSamples int
	MinRun          int64
	MaxRun          int64
	Median          int64
	Mean            int64
	Runs            []int64
}
