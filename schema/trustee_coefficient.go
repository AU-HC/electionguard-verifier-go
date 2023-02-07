package schema

import "electionguard-verifier-go/crypto"

type TrusteeCoefficient struct {
	PublicKey int          // The public key generated from secret coefficient TODO: Change type
	Proof     crypto.Proof // Schnorr proof for the private key TODO: Change type
}
