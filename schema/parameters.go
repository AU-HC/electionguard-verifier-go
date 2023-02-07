package schema

import "math/big"

// Parameters necessary to form the election.
type Parameters struct {
	date          string  // The date of the election
	location      string  // The location of the election
	numOfTrustees int     // The number of guardians (namely 'n')
	threshold     int     // The threshold number to decrypt the aggregate ballot 'k'
	prime         big.Int // The prime modulus of the group, which is used for encryption

	// Need generator
	// generator
}
