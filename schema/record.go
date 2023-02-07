package schema

import "math/big"

type Record struct {
	parameters Parameters // The parameters for the election
	baseHash   big.Int    // The base hash 'Q'

}
