package test

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"testing"
)

func TestHMACParameterHash(t *testing.T) {
	ver := append([]byte("v2.0.0"), make([]byte, 27)...)

	got := crypto.HMAC(*schema.MakeBigIntFromByteArray(ver), 0x00, schema.BigInt{})

	wanted := new(schema.BigInt)
	wanted.SetString("", 10)

	hashIsIncorrect := !got.Compare(wanted)
	if hashIsIncorrect {
		t.Error(formatBigIntErrorString(got, wanted))
	}
}
