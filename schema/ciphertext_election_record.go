package schema

type CiphertextElectionRecord struct {
	NumberOfGuardians      int          `json:"number_of_guardians"`
	Quorum                 int          `json:"quorum"`
	ElgamalPublicKey       BigInt       `json:"elgamal_public_key"`
	CommitmentHash         BigInt       `json:"commitment_hash"`
	ManifestHash           BigInt       `json:"manifest_hash"`
	CryptoBaseHash         BigInt       `json:"crypto_base_hash"`
	CryptoExtendedBaseHash BigInt       `json:"crypto_extended_base_hash"`
	ExtendedData           ExtendedData `json:"extended_data"`
	Configuration          struct {
		AllowOvervotes bool `json:"allow_overvotes"`
		MaxVotes       int  `json:"max_votes"`
	} `json:"configuration"`
}
