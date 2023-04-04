package utility

import (
	"flag"
)

type ApplicationArguments struct {
	LoggingLevel              LoggingLevel
	ElectionArtifactsPath     string
	OutputPath                string
	ConcurrentSteps           bool
	AmountBenchmarkingSamples int
}

func InitApplicationArguments() ApplicationArguments {
	// Creating struct with empty arguments
	arguments := ApplicationArguments{}

	// Getting arguments from flags
	loggingLevelIntPtr := flag.Int("v", 0, "Logging level: 0 = no logging, 1 = info and higher, 2 = debug and higher")
	flag.StringVar(&arguments.ElectionArtifactsPath, "p", "data/hamilton-general/election_record", "Path to election record")
	flag.StringVar(&arguments.OutputPath, "o", "", "File which to output verification result")
	flag.BoolVar(&arguments.ConcurrentSteps, "c", true, "Decides if the verifier should run the verification steps concurrent")
	flag.IntVar(&arguments.AmountBenchmarkingSamples, "b", 0, "Decides if the verifier should be benchmarked, and the amount of samples")

	// Parsing flags
	flag.Parse()
	arguments.LoggingLevel = intToLoggingLevel(*loggingLevelIntPtr)

	return arguments
}

func intToLoggingLevel(level int) LoggingLevel {
	switch level {
	case 1:
		return LogInfo
	case 2:
		return LogDebug
	default:
		return LogNone
	}
}
