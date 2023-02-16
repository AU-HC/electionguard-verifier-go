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
// should however consider calculating ((a mod n)+(b mod n)) mod n, this returns a new BigInt
func (b *BigInt) ModAddition(a *BigInt, n *BigInt) *BigInt {
	// Start by adding the two numbers
	resultOfAddition := b.Add(&b.Int, &a.Int)

	// Calculate (b+a) mod n
	resultOfMod := b.Mod(resultOfAddition, &n.Int)

	// Return new BigInt
	var bigIntResult BigInt
	bigIntResult.Int = *resultOfMod
	return &bigIntResult
}

// ModMul is used to calculate (a*b) mod n <=> ((a mod n) · (b mod n)) mod n
// should however consider calculating ((a mod n)+(b mod n)) mod n
func ModMul(a *BigInt, b, n *BigInt) *BigInt {
	// Start by taking the mod of both numbers
	modOfA := a.Mod(&a.Int, &n.Int)
	modOfB := b.Mod(&b.Int, &n.Int)

	// Multiply the two numbers mod n
	resultOfMultiplication := modOfA.Mul(modOfA, modOfB)
	resultOfMod := resultOfMultiplication.Mod(resultOfMultiplication, &n.Int)

	// Return new BigInt
	var bigIntResult BigInt
	bigIntResult.Int = *resultOfMod
	return &bigIntResult
}

func ModExp(base, exponent, modulus *BigInt) *BigInt {
	var result big.Int
	result.SetString("1", 10)
	var zero big.Int

	// Checking if exponent > 0
	for (zero.Cmp(&exponent.Int)) < 0 {
		// If y is odd, multiply base with result
		if exponent.Bit(0) == 1 {
			result.Mul(&result, &base.Int)
		}

		// Exponent is even now, so divide by 2 and base^2
		exponent.Rsh(&exponent.Int, 1)
		base.Mul(&base.Int, &base.Int)
	}
	// Result mod p
	resultOfMod := result.Mod(&result, &modulus.Int)
	var bigIntResult BigInt
	bigIntResult.Int = *resultOfMod

	return &bigIntResult
}
