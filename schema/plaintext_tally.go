package schema

type PlaintextTally struct {
	ObjectId string                  `json:"object_id"`
	Contests map[string]ContestTally `json:"contests"`
}

type ContestTally struct {
	ObjectId   string                    `json:"object_id"`
	Selections map[string]SelectionTally `json:"selections"`
}

type SelectionTally struct {
	ObjectId string     `json:"object_id"`
	Tally    int        `json:"tally"`
	Value    string     `json:"value"`
	Message  Ciphertext `json:"message"`
	Shares   []struct {
		ObjectId       string      `json:"object_id"`
		GuardianId     string      `json:"guardian_id"`
		Share          string      `json:"share"`
		Proof          CpProof     `json:"proof"`
		RecoveredParts interface{} `json:"recovered_parts"`
	} `json:"shares"`
}

type CpProof struct {
	Pad       BigInt `json:"pad"`
	Data      BigInt `json:"data"`
	Challenge BigInt `json:"challenge"`
	Response  BigInt `json:"response"`
	Usage     string `json:"usage"`
}
