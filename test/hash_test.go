package test

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"testing"
)

func TestHashThreeInt(t *testing.T) {
	got := crypto.HashElems(1, 2, 3)
	wanted := new(schema.BigInt)
	wanted.SetString("101860255573162687554529317338421470715872259126123982930947143077424977508731", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func TestHashOneInt(t *testing.T) {
	got := crypto.HashElems(1)
	wanted := new(schema.BigInt)
	wanted.Int.SetString("55842377753952778025173527915631301100693874962723145527793793795385961650435", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func formatBigIntErrorString(got, wanted *schema.BigInt) string {
	return "Got: " + got.String() + "\nWanted: " + wanted.String()
}
