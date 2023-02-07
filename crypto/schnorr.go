package crypto

type Proof struct {
	Commitment int // g^r which acts as a Commitment to the random number 'r' TODO: Change type
	Challenge  int // The challenge c, which is a hash of relevant parameters TODO: Change type
	Response   int // The response z = r+cx TODO: Change type and comment
}

// Check that the sender has possession over the private key which corresponds to the public key.
func (p *Proof) Check() {
	// TODO: Implement
}
