package utility

type LoggingLevel int64

const (
	LogNone LoggingLevel = iota
	LogInfo
	LogDebug
)

const SampleDataDir = "data/hamilton-general/election_record"
