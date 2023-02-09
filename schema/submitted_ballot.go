package schema

type SubmittedBallots struct {
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
				ProofZeroPad       string `json:"proof_zero_pad"` // TODO: Change to big.Int?
				ProofZeroData      string `json:"proof_zero_data"`
				ProofOnePad        string `json:"proof_one_pad"`
				ProofOneData       string `json:"proof_one_data"`
				ProofZeroChallenge string `json:"proof_zero_challenge"`
				ProofOneChallenge  string `json:"proof_one_challenge"`
				Challenge          string `json:"challenge"`
				ProofZeroResponse  string `json:"proof_zero_response"`
				ProofOneResponse   string `json:"proof_one_response"`
				Usage              string `json:"usage"`
			} `json:"proof"`
			ExtendedData interface{} `json:"extended_data"`
		} `json:"ballot_selections"`
		CiphertextAccumulation struct {
			Pad  string `json:"pad"`
			Data string `json:"data"`
		} `json:"ciphertext_accumulation"`
		CryptoHash string      `json:"crypto_hash"`
		Nonce      interface{} `json:"nonce"`
		Proof      struct {
			Pad       string `json:"pad"`
			Data      string `json:"data"`
			Challenge string `json:"challenge"`
			Response  string `json:"response"`
			Constant  int    `json:"constant"`
			Usage     string `json:"usage"`
		} `json:"proof"`
		ExtendedData interface{} `json:"extended_data"`
	} `json:"contests"`
	Code       string      `json:"code"`
	Timestamp  int         `json:"timestamp"`
	CryptoHash string      `json:"crypto_hash"`
	Nonce      interface{} `json:"nonce"`
	State      int         `json:"state"`
}
