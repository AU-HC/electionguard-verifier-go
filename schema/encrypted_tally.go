package schema

type EncryptedTally struct {
	ObjectId string                            `json:"object_id"`
	Contests map[string]EncryptedContestsTally `json:"contests"`
}

type EncryptedContestsTally struct {
	ObjectId        string                             `json:"object_id"`
	SequenceOrder   int                                `json:"sequence_order"`
	DescriptionHash BigInt                             `json:"description_hash"`
	Selections      map[string]EncryptedSelectionTally `json:"selections"`
}

type EncryptedSelectionTally struct {
	ObjectId        string     `json:"object_id"`
	SequenceOrder   int        `json:"sequence_order"`
	DescriptionHash BigInt     `json:"description_hash"`
	Ciphertext      Ciphertext `json:"ciphertext"`
}

type Ciphertext struct {
	Pad  BigInt `json:"pad"`
	Data BigInt `json:"data"`
}
