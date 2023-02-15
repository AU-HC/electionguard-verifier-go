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

func TestHashString(t *testing.T) {
	got := crypto.HashElems("string")
	wanted := new(schema.BigInt)
	wanted.Int.SetString("90926586383276802466644404271371545801279822268723715256179656356152275330028", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func TestHashNil(t *testing.T) {
	got := crypto.HashElems(nil)
	wanted := new(schema.BigInt)
	wanted.Int.SetString("34190542803364976751518993874547876265610841613775469638084026275531073571566", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func TestHashBigInt(t *testing.T) {
	toBeHashed := new(schema.BigInt)
	toBeHashed.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF43", 16)
	got := crypto.HashElems(*toBeHashed) // Import, cannot pass pointer, must dereference

	wanted := new(schema.BigInt)
	wanted.Int.SetString("81214711932601036140639252629910130067642288172673877624082744215981851250268", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func TestHashIntSlice(t *testing.T) {
	toBeHashed := [3]string{"hej", "med", "dig"}
	got := crypto.HashElems(toBeHashed)

	wanted := new(schema.BigInt)
	wanted.Int.SetString("69616850468205167024114498771676296544077351555488666079379595094599566413508", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}

func formatBigIntErrorString(got, wanted *schema.BigInt) string {
	return "Got: " + got.String() + "\nWanted: " + wanted.String()
}
