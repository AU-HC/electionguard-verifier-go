package schema

type SpoiledBallot struct {
	ObjectId string                      `json:"object_id"`
	Contests map[string]DecryptedContest `json:"contests"`
}

type DecryptedContest struct {
	ObjectId    string                        `json:"object_id"`
	Selections  map[string]DecryptedSelection `json:"selections"`
	ContestData ContestData                   `json:"contest_data"`
}

type ContestData struct {
	ObjectId      string           `json:"object_id"`
	DecryptedData BigInt           `json:"data"`
	Ciphertext    ExtendedData     `json:"ciphertext_extended_data"`
	Shares        []SelectionShare `json:"shares"`
}

type DecryptedSelection struct {
	ObjectId string           `json:"object_id"`
	Tally    int              `json:"tally"`
	Value    BigInt           `json:"value"`
	Message  Ciphertext       `json:"message"`
	Shares   []SelectionShare `json:"shares"`
}
