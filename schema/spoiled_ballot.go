package schema

// TODO: Check this whole thing :)))))))

type SpoiledBallot struct {
	ObjectId string                      `json:"object_id"`
	Contests map[string]DecryptedContest `json:"contests"`
}

type DecryptedContest struct {
	ObjectId    string                        `json:"object_id"`
	Selections  map[string]DecryptedSelection `json:"selections"`
	ContestData struct {
		ObjectId string `json:"object_id"`
		Data     string `json:"data"`
	}
}

type DecryptedSelection struct {
	ObjectId string `json:"object_id"`
	Tally    int    `json:"tally"`
	Value    string `json:"value"`
	Message  struct {
		Pad  string `json:"pad"`
		Data string `json:"data"`
	} `json:"message"`
	Shares []struct {
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
