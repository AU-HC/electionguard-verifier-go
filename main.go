package main

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/utility"
	"fmt"
	"runtime"
)

func main() {
	// Fetching flags and amount of logical cores for the CPU
	applicationArguments := utility.InitApplicationArguments()
	path := applicationArguments.ElectionArtifactsPath
	amountOfLogicalCores := runtime.NumCPU()

	// Fetching logging level and creating logger
	logger := utility.ConfigureLogger(applicationArguments.LoggingLevel)
	outputStrategy := core.MakeOutputStrategy(applicationArguments.OutputPath)
	verifyStrategy := core.MakeVerifyStrategy(applicationArguments.ConcurrentSteps, amountOfLogicalCores)
	samples := applicationArguments.AmountBenchmarkingSamples

	// Create verifier, set strategies, and verify election data
	verifier := *core.MakeVerifier(logger)
	verifier.SetOutputStrategy(outputStrategy)
	verifier.SetVerifyStrategy(verifyStrategy)

	// Should probably refactor
	if samples == 0 {
		electionIsValid := verifier.Verify(path)

		// Result of verification of the election data
		if electionIsValid {
			fmt.Println("Election is valid")
		} else {
			fmt.Println("Election is invalid")
		}
	} else {
		// verifier.Benchmark(path, samples)
		verifier.BenchmarkDeserialization(samples)
	}
}
