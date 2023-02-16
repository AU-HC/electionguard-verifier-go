package test

import (
	"electionguard-verifier-go/schema"
	"testing"
)

// TODO: Should change to larger integers
func TestModMul(t *testing.T) {
	a := new(schema.BigInt)
	a.SetString("3", 10)
	b := new(schema.BigInt)
	b.SetString("4", 10)
	mod := new(schema.BigInt)
	mod.SetString("5", 10)

	got := schema.ModMul(a, b, mod)

	wanted := new(schema.BigInt)
	wanted.SetString("2", 10)

	gotIsDifferentFromWanted := !got.Compare(wanted)
	if gotIsDifferentFromWanted {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func TestModExp(t *testing.T) {
	base := new(schema.BigInt)
	base.SetString("2", 10)
	exp := new(schema.BigInt)
	exp.SetString("5", 10)
	mod := new(schema.BigInt)
	mod.SetString("13", 10)

	got := schema.ModExp(base, exp, mod)

	wanted := new(schema.BigInt)
	wanted.SetString("6", 10)

	gotIsDifferentFromWanted := !got.Compare(wanted)
	if gotIsDifferentFromWanted {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}
