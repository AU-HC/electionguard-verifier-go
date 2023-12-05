package schema

type Context struct {
	ManifestHash           BigInt        `json:"manifest_hash"`
	CommitmentHash         BigInt        `json:"commitment_hash"`
	CryptoBaseHash         BigInt        `json:"crypto_base_hash"`
	CryptoExtendedBaseHash BigInt        `json:"crypto_extended_base_hash"`
	ElgamalPublicKey       BigInt        `json:"elgamal_public_key"`
	NumberOfGuardians      int           `json:"number_of_guardians"`
	Quorum                 int           `json:"quorum"`
	Configuration          Configuration `json:"configuration"`
}

type Configuration struct {
	AllowOvervotes bool `json:"allow_overvotes"`
	MaxVotes       int  `json:"max_votes"`
}
