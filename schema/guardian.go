package schema

type Guardian struct {
	GuardianId          string         `json:"guardian_id"`
	SequenceOrder       int            `json:"sequence_order"`
	ElectionPublicKey   BigInt         `json:"election_public_key"`
	ElectionCommitments []BigInt       `json:"election_commitments"`
	ElectionProofs      []SchnorrProof `json:"election_proofs"`
}

type SchnorrProof struct {
	PublicKey  BigInt `json:"public_key"`
	Commitment BigInt `json:"commitment"`
	Challenge  BigInt `json:"challenge"`
	Response   BigInt `json:"response"`
	Usage      string `json:"usage"`
}
