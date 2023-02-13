package schema

type SubmittedBallot struct {
	ObjectId     string `json:"object_id"`
	StyleId      string `json:"style_id"`
	ManifestHash string `json:"manifest_hash"`
	CodeSeed     string `json:"code_seed"`
	Contests     []struct {
		ObjectId         string `json:"object_id"`
		SequenceOrder    int    `json:"sequence_order"`
		DescriptionHash  string `json:"description_hash"`
		BallotSelections []struct {
			ObjectId               string      `json:"object_id"`
			SequenceOrder          int         `json:"sequence_order"`
			DescriptionHash        string      `json:"description_hash"`
			Ciphertext             Ciphertext  `json:"ciphertext"`
			CryptoHash             string      `json:"crypto_hash"`
			IsPlaceholderSelection bool        `json:"is_placeholder_selection"`
			Nonce                  interface{} `json:"nonce"`
			Proof                  struct {
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
			} `json:"proof"`
			ExtendedData interface{} `json:"extended_data"`
		} `json:"ballot_selections"`
		CiphertextAccumulation struct {
			Pad  BigInt `json:"pad"`
			Data BigInt `json:"data"`
		} `json:"ciphertext_accumulation"`
		CryptoHash   string        `json:"crypto_hash"`
		Nonce        interface{}   `json:"nonce"`
		Proof        ConstantProof `json:"proof"`
		ExtendedData interface{}   `json:"extended_data"`
	} `json:"contests"`
	Code       string      `json:"code"`
	Timestamp  int         `json:"timestamp"`
	CryptoHash string      `json:"crypto_hash"`
	Nonce      interface{} `json:"nonce"`
	State      int         `json:"state"`
}

type ConstantProof struct {
	Pad       BigInt `json:"pad"`
	Data      BigInt `json:"data"`
	Challenge BigInt `json:"challenge"`
	Response  BigInt `json:"response"`
	Constant  int    `json:"constant"`
	Usage     string `json:"usage"`
}
