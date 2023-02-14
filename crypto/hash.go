package crypto

import (
	"crypto/sha256"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

type SHA256 struct {
	hash [32]byte
	q    schema.BigInt
}

func MakeSHA256() *SHA256 {
	return &SHA256{q: utility.MakeCorrectElectionConstants().Q}
}

func (s *SHA256) update(data string) {
	// Hashing the new data string and converting it to []byte from [32]byte
	var hash32 = sha256.Sum256([]byte(data))
	hash := hash32[:]

	// Creating big.Int for the modular addition
	intCurrentHash := new(big.Int).SetBytes(hash)
	intPreviousHash := new(big.Int).SetBytes(s.hash[:])
	currentHashBigInt := schema.BigInt{Int: *intCurrentHash}
	previousHashBigInt := schema.BigInt{Int: *intPreviousHash}

	// Creating new bigInt
	modularAddition := currentHashBigInt.ModAddition(&previousHashBigInt, &s.q)
	modularAddition.FillBytes(s.hash[:])
}

var stringType = reflect.TypeOf("")
var intType = reflect.TypeOf(1)
var bigIntType = reflect.TypeOf(schema.BigInt{})

// var intSliceType = reflect.TypeOf(([]int)(nil))
// var submittedBallotSliceType = reflect.TypeOf(([]schema.SubmittedBallot)(nil))

func HashElems(a ...interface{}) schema.BigInt {
	// StringBuilder
	h := MakeSHA256()

	for _, x := range a {
		var hashMe string

		switch reflect.TypeOf(x) {
		case intType:
			xInt, _ := x.(int)
			hashMe = strconv.Itoa(xInt)
		case stringType:
			hashMe, _ = x.(string) // type cast (strings are already utf8-encoded in Golang)
		case bigIntType:
			bigInt := x.(schema.BigInt).Int
			hex := fmt.Sprintf("%X", &bigInt)    // Convert big.Int to hex
			hashMe = addLeadingZeroIfNeeded(hex) // add leading zero if amount of digits is odd
		}

		h.update(hashMe + "|")
	}

	var result big.Int
	result.SetBytes(h.hash[:])
	return schema.BigInt{Int: result}
}

func addLeadingZeroIfNeeded(hex string) string {
	stringLengthIsOdd := len(hex)%2 != 0
	if stringLengthIsOdd {
		return hex
	}

	return "0" + hex
}
