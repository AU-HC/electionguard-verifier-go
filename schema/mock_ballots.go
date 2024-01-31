package schema

type MockBallot struct {
	ObjectId string              `json:"ObjectId"`
	Code     string              `json:"Code"`
	Contests []MockBallotContest `json:"Contests"`
}

type MockBallotContest struct {
	ObjectId         string                `json:"ObjectId"`
	BallotSelections []MockBallotSelection `json:"BallotSelections"`
	Proof            MockRangeProof        `json:"Proof"`
}

type MockBallotSelection struct {
	ObjectId   string         `json:"ObjectId"`
	Ciphertext Ciphertext     `json:"Ciphertext"`
	Proof      MockRangeProof `json:"Proof"`
}

type MockRangeProof struct {
	Challenge  BigInt               `json:"Challenge"`
	Proofs     []ChaumPedersenProof `json:"Proofs"`
	RangeLimit int                  `json:"RangeLimit"`
}
