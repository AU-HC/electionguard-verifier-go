package core

import (
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
	"math/big"
)

func isValidResidue(a schema.BigInt) bool {
	// Checking the value is in range
	cons := utility.MakeCorrectElectionConstants()
	p := cons.P
	q := cons.Q
	zero := big.NewInt(0)
	var one schema.BigInt
	one.SetString("1", 10)

	valueIsAboveOrEqualToZero := zero.Cmp(&a.Int) <= 0
	valueIsSmallerThanP := p.Cmp(&a.Int) == 1
	validResidue := (powP(&a, &q)).Compare(&one)

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP && validResidue
}

func isInRange(a schema.BigInt) bool {
	q := utility.MakeCorrectElectionConstants().Q.Int
	zero := big.NewInt(0)

	valueIsAboveOrEqualToZero := zero.Cmp(&a.Int) <= 0
	valueIsSmallerThanP := q.Cmp(&a.Int) > 0

	return valueIsAboveOrEqualToZero && valueIsSmallerThanP
}

func powP(b, e *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := utility.MakeCorrectElectionConstants().P.Int

	result.Exp(&b.Int, &e.Int, &p)

	return &result
}

func mulP(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	p := utility.MakeCorrectElectionConstants().P.Int

	modOfA := a.Mod(&a.Int, &p)
	modOfB := b.Mod(&b.Int, &p)

	// Multiply the two numbers mod n
	resultOfMultiplication := modOfA.Mul(modOfA, modOfB)
	result.Mod(resultOfMultiplication, &p)

	return &result
}

func addQ(a, b *schema.BigInt) *schema.BigInt {
	var result schema.BigInt
	q := utility.MakeCorrectElectionConstants().Q.Int

	result.Add(&b.Int, &a.Int)
	result.Mod(&result.Int, &q)

	return &result
}
