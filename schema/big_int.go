package schema

import (
	"fmt"
	"math/big"
	"strings"
)

type BigInt struct {
	big.Int
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	// If string is null, return nil
	if string(p) == "null" {
		return nil
	}

	s := convertByteArrayToStringAndTrim(p)

	// Declare new bigInt and return error if SetString fails
	z := big.Int{}
	_, ok := z.SetString(s, 16) // TODO: Check base
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z

	// Return no error
	return nil
}

func convertByteArrayToStringAndTrim(p []byte) string {
	s := string(p)
	s = strings.TrimRight(s, "\"")
	s = strings.TrimLeft(s, "\"")

	return s
}

// Compare Method to make true/false comparisons easier
func (b *BigInt) Compare(a *BigInt) bool {
	isEqual := b.Cmp(&a.Int)
	if isEqual == 0 {
		return true
	}

	return false
}

// ModAddition is used to calculate (b+a) mod n
// should however consider calculating ((a mod n)+(b mod n)) mod n
func (b *BigInt) ModAddition(a *BigInt, n *BigInt) *BigInt {
	// Start by adding the two numbers
	resultOfAddition := b.Add(&b.Int, &a.Int)

	// Calculate (b+a) mod n
	resultOfMod := b.Mod(resultOfAddition, &n.Int)

	// Return the BigInt
	b.Int = *resultOfMod
	return b
}

// ModMul is used to calculate (b+a) mod n
// should however consider calculating ((a mod n)+(b mod n)) mod n
func (b *BigInt) ModMul(a *BigInt, n *BigInt) *BigInt {
	return b
}

func ModularExponentiation(base, exponent, modulus *BigInt) *BigInt {
	return base
}
