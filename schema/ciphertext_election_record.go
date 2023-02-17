package schema

type CiphertextElectionRecord struct {
	NumberOfGuardians      int         `json:"number_of_guardians"`
	Quorum                 int         `json:"quorum"`
	ElgamalPublicKey       BigInt      `json:"elgamal_public_key"` // TODO: Consider changing type
	CommitmentHash         []byte      `json:"commitment_hash"`
	ManifestHash           []byte      `json:"manifest_hash"`
	CryptoBaseHash         []byte      `json:"crypto_base_hash"`
	CryptoExtendedBaseHash []byte      `json:"crypto_extended_base_hash"`
	ExtendedData           interface{} `json:"extended_data"` // TODO: Consider changing type
	Configuration          struct {
		AllowOvervotes bool `json:"allow_overvotes"`
		MaxVotes       int  `json:"max_votes"`
	} `json:"configuration"`
}
