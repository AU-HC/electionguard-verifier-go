package schema

type SubmittedBallot struct {
	ObjectId string          `json:"object_id"`
	Code     string          `json:"code"`
	Contests []BallotContest `json:"contests"`
}

type BallotContest struct {
	ObjectId               string            `json:"object_id"`
	CryptoHash             BigInt            `json:"crypto_hash"`
	BallotSelections       []BallotSelection `json:"ballot_selections"`
	Proof                  RangeProof        `json:"proof"`
	CiphertextAccumulation Ciphertext        `json:"ciphertext_accumulation"`
}

type BallotSelection struct {
	ObjectId     string           `json:"object_id"`
	CryptoHash   BigInt           `json:"crypto_hash"`
	Proof        DisjunctiveProof `json:"proof"`
	Ciphertext   Ciphertext       `json:"ciphertext"`
	ExtendedData ExtendedData     `json:"extended_data"`
}

type DisjunctiveProof struct {
	Challenge          BigInt `json:"challenge"`
	ProofZeroPad       BigInt `json:"proof_zero_pad"`
	ProofZeroData      BigInt `json:"proof_zero_data"`
	ProofZeroChallenge BigInt `json:"proof_zero_challenge"`
	ProofZeroResponse  BigInt `json:"proof_zero_response"`
	ProofOnePad        BigInt `json:"proof_one_pad"`
	ProofOneData       BigInt `json:"proof_one_data"`
	ProofOneChallenge  BigInt `json:"proof_one_challenge"`
	ProofOneResponse   BigInt `json:"proof_one_response"`
}

type RangeProof struct {
	Challenge  BigInt          `json:"challenge"`
	Proofs     [][]interface{} `json:"proofs"` // TODO?
	RangeLimit int             `json:"range_limit"`
}
