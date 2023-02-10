package utility

type LoggingLevel int64

const (
	LogDebug LoggingLevel = iota
	LogInfo
	LogNone
)

const SAMPLE_DATA_DIR = "data/hamilton-general/election_record"
