package utility

type LoggingLevel int64

const (
	LogDebug LoggingLevel = iota
	LogInfo
	LogNone
)

const SampleDataDir = "data/hamilton-general/election_record"
