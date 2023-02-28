package schema

type SubmittedBallot struct {
	ObjectId     string          `json:"object_id"`
	StyleId      string          `json:"style_id"`
	ManifestHash BigInt          `json:"manifest_hash"`
	Seed         BigInt          `json:"code_seed"`
	Contests     []BallotContest `json:"contests"`
	Code         BigInt          `json:"code"`
	Timestamp    int             `json:"timestamp"`
	CryptoHash   BigInt          `json:"crypto_hash"`
	Nonce        interface{}     `json:"nonce"`
	State        int             `json:"state"`
}

type BallotContest struct {
	ObjectId               string            `json:"object_id"`
	SequenceOrder          int               `json:"sequence_order"`
	DescriptionHash        BigInt            `json:"description_hash"`
	BallotSelections       []BallotSelection `json:"ballot_selections"`
	CiphertextAccumulation Ciphertext        `json:"ciphertext_accumulation"`
	CryptoHash             BigInt            `json:"crypto_hash"`
	Nonce                  interface{}       `json:"nonce"`
	Proof                  ConstantProof     `json:"proof"`
	ExtendedData           ExtendedData      `json:"extended_data"`
}

type BallotSelection struct {
	ObjectId               string           `json:"object_id"`
	SequenceOrder          int              `json:"sequence_order"`
	DescriptionHash        BigInt           `json:"description_hash"`
	Ciphertext             Ciphertext       `json:"ciphertext"`
	CryptoHash             BigInt           `json:"crypto_hash"`
	IsPlaceholderSelection bool             `json:"is_placeholder_selection"`
	Nonce                  interface{}      `json:"nonce"`
	Proof                  DisjunctiveProof `json:"proof"`
}

type DisjunctiveProof struct {
	ProofZeroPad       BigInt `json:"proof_zero_pad"`
	ProofZeroData      BigInt `json:"proof_zero_data"`
	ProofOnePad        BigInt `json:"proof_one_pad"`
	ProofOneData       BigInt `json:"proof_one_data"`
	ProofZeroChallenge BigInt `json:"proof_zero_challenge"`
	ProofOneChallenge  BigInt `json:"proof_one_challenge"`
	Challenge          BigInt `json:"challenge"`
	ProofZeroResponse  BigInt `json:"proof_zero_response"`
	ProofOneResponse   BigInt `json:"proof_one_response"`
	Usage              string `json:"usage"`
}

type ConstantProof struct {
	Pad       BigInt `json:"pad"`
	Data      BigInt `json:"data"`
	Challenge BigInt `json:"challenge"`
	Response  BigInt `json:"response"`
	Constant  int    `json:"constant"`
	Usage     string `json:"usage"`
}
