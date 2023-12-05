package core

import (
	"electionguard-verifier-go/schema"
	"math/big"
)

func (v *Verifier) isValidResidue(a schema.BigInt) bool {
	// Checking the value is in range
	p := &v.constants.P.Int
	q := v.constants.Q
	zero := schema.IntToBigInt(0)
	one := schema.IntToBigInt(1)

	valueIsAboveOrEqualToZero := zero.Cmp(&a.Int) <= 0
	valueIsSmallerThanP := p.Cmp(&a.Int) == 1
	valueIsInRange := valueIsSmallerThanP && valueIsAboveOrEqualToZero // a is in [0, P)

	validResidue := v.powP(&a, q).Compare(one) // a^q mod p == 1

	return valueIsInRange && validResidue
}

func (v *Verifier) isInRange(a schema.BigInt) bool {
	q := &v.constants.Q.Int
	zero := big.NewInt(0)

	valueIsAboveOrEqualToZero := zero.Cmp(&a.Int) <= 0
	valueIsSmallerThanP := q.Cmp(&a.Int) > 0

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP
}

// sub returns a-b
func sub(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	result.Sub(&a.Int, &b.Int)

	return &result
}

func mul(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	result.Mul(&a.Int, &b.Int)

	return &result
}

func modP(a *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := GetP()

	result.Mod(&a.Int, &p.Int)
	return &result
}

func modQ(a *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	q := GetQ()

	result.Mod(&a.Int, &q.Int)
	return &result
}

func (v *Verifier) powP(b, e *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := &v.constants.P.Int

	result.Exp(&b.Int, &e.Int, p)

	return &result
}

func (v *Verifier) mulP(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := &v.constants.P.Int

	modOfA := a.Mod(&a.Int, p)
	modOfB := b.Mod(&b.Int, p)

	// Multiply the two numbers mod n
	result.Mul(modOfA, modOfB)
	result.Mod(&result.Int, p)

	return &result
}

func (v *Verifier) addQ(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	q := &v.constants.Q.Int

	result.Add(&b.Int, &a.Int)
	result.Mod(&result.Int, q)

	return &result
}

func (v *Verifier) addP(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := &v.constants.P.Int

	result.Add(&b.Int, &a.Int)
	result.Mod(&result.Int, p)

	return &result
}

func (v *Verifier) subQ(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	q := &v.constants.Q.Int

	result.Sub(&a.Int, &b.Int)
	result.Mod(&result.Int, q)

	return &result
}

func (v *Verifier) invP(a *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := &v.constants.P.Int

	result.ModInverse(&a.Int, p)
	return &result
}
