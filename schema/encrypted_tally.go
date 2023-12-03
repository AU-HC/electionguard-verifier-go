package schema

type EncryptedTally struct {
	TallyID  string                            `json:"tally_id"`
	Name     string                            `json:"name"`
	Contests map[string]EncryptedContestsTally `json:"contests"`
}

type EncryptedContestsTally struct {
	ObjectId   string                             `json:"object_id"`
	Selections map[string]EncryptedSelectionTally `json:"selections"`
}

type EncryptedSelectionTally struct {
	ObjectId   string     `json:"object_id"`
	Ciphertext Ciphertext `json:"ciphertext"`
}

type Ciphertext struct {
	Pad  BigInt `json:"pad"`
	Data BigInt `json:"data"`
}
