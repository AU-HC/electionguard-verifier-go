package utility

import (
	"flag"
)

type ApplicationArguments struct {
	LoggingLevel          LoggingLevel
	ElectionArtifactsPath string
	OutputPath            string
	UseMultipleThreads    bool
}

func InitApplicationArguments() ApplicationArguments {
	// Creating struct with empty arguments
	arguments := ApplicationArguments{}

	// Getting arguments from flags
	loggingLevelIntPtr := flag.Int("v", 0, "Logging level: 0 = no logging, 1 = info and higher, 2 = debug and higher")
	flag.StringVar(&arguments.ElectionArtifactsPath, "p", "data/hamilton-general/election_record", "Path to election record")
	flag.StringVar(&arguments.OutputPath, "o", "", "File which to output verification result")
	flag.BoolVar(&arguments.UseMultipleThreads, "c", true, "Decides if the verifier should run the verification steps concurrent")

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
