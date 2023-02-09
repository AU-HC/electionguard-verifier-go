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
		ObjectId   string `json:"object_id"`
		GuardianId string `json:"guardian_id"`
		Share      string `json:"share"`
		Proof      struct {
			Pad       string `json:"pad"`
			Data      string `json:"data"`
			Challenge string `json:"challenge"`
			Response  string `json:"response"`
			Usage     string `json:"usage"`
		} `json:"proof"`
		RecoveredParts interface{} `json:"recovered_parts"`
	} `json:"shares"`
}
