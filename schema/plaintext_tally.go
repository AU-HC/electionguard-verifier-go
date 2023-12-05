package schema

type PlaintextTally struct {
	TallyID  string                      `json:"tally_id"`
	Name     string                      `json:"name"`
	Contests map[string]DecryptedContest `json:"contests"`
}
