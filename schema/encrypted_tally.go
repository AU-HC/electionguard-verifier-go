package schema

type EncryptedTally struct {
	ObjectId string              `json:"object_id"`
	Contests map[string]Contests `json:"contests"`
}

type Contests struct {
	ObjectId        string                `json:"object_id"`
	SequenceOrder   int                   `json:"sequence_order"`
	DescriptionHash string                `json:"description_hash"`
	Selections      map[string]Selections `json:"selections"`
}

type Selections struct {
	ObjectId        string     `json:"object_id"`
	SequenceOrder   int        `json:"sequence_order"`
	DescriptionHash string     `json:"description_hash"`
	Ciphertext      Ciphertext `json:"ciphertext"`
}

type Ciphertext struct {
	Pad  string `json:"pad"`
	Data string `json:"data"`
}
