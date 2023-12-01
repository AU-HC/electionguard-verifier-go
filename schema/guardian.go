package schema

type Guardian struct {
	ObjectID                string         `json:"object_id"`
	GuardianId              string         `json:"guardian_id"`
	SequenceOrder           int            `json:"sequence_order"`
	Key                     BigInt         `json:"key"`
	CoefficientsCommitments []BigInt       `json:"coefficient_commitments"`
	CoefficientsProofs      []SchnorrProof `json:"coefficient_proofs"`
}

type SchnorrProof struct {
	PublicKey  BigInt `json:"public_key"`
	Commitment BigInt `json:"commitment"`
	Challenge  BigInt `json:"challenge"`
	Response   BigInt `json:"response"`
	Usage      int    `json:"usage"`
}
