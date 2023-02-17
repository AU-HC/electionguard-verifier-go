package core

import (
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
)

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
