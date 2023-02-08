package schema

type Guardian struct {
	GuardianId          string   `json:"guardian_id"`
	SequenceOrder       int      `json:"sequence_order"`
	ElectionPublicKey   string   `json:"election_public_key"`
	ElectionCommitments []string `json:"election_commitments"`
	ElectionProofs      []struct {
		PublicKey  string `json:"public_key"`
		Commitment string `json:"commitment"`
		Challenge  string `json:"challenge"`
		Response   string `json:"response"`
		Usage      string `json:"usage"`
	} `json:"election_proofs"`
}
