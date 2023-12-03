package schema

type SpoiledBallot struct {
	BallotID string                      `json:"ballot_id"`
	TallyID  string                      `json:"tally_id"`
	Name     string                      `json:"name"`
	Contests map[string]DecryptedContest `json:"contests"`
}

type DecryptedContest struct {
	ObjectId   string                        `json:"object_id"`
	Selections map[string]DecryptedSelection `json:"selections"`
}

type DecryptedSelection struct {
	ObjectId string             `json:"object_id"`
	Tally    int                `json:"tally"`
	Value    BigInt             `json:"value"`
	Proof    ChaumPedersenProof `json:"proof"`
}

type ChaumPedersenProof struct {
	Pad       BigInt `json:"pad"`
	Data      BigInt `json:"data"`
	Challenge BigInt `json:"challenge"`
	Response  BigInt `json:"response"`
}
