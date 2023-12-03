package schema

type PlaintextTally struct {
	TallyID  string                  `json:"tally_id"`
	Name     string                  `json:"name"`
	Contests map[string]ContestTally `json:"contests"`
}

type ContestTally struct {
	ObjectId   string                    `json:"object_id"`
	Selections map[string]SelectionTally `json:"selections"`
}

type SelectionTally struct {
	ObjectId string             `json:"object_id"`
	Tally    int                `json:"tally"`
	Value    BigInt             `json:"value"`
	Proof    ChaumPedersenProof `json:"proof"`
}
