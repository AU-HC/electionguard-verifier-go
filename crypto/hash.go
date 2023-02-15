package crypto

import (
	"bytes"
	"crypto/sha256"
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

var nilType = reflect.TypeOf(nil)
var stringType = reflect.TypeOf("")
var intType = reflect.TypeOf(1)
var intSliceType = reflect.TypeOf(([]int)(nil))
var bigIntType = reflect.TypeOf(schema.BigInt{})
var bigIntSliceType = reflect.TypeOf(([]schema.BigInt)(nil))

// var submittedBallotSliceType = reflect.TypeOf(([]schema.SubmittedBallot)(nil))

type SHA256 struct {
	toHash bytes.Buffer
	q      schema.BigInt // TODO: Probably optimize the way Q is handled
}

func MakeSHA256() *SHA256 {
	return &SHA256{q: utility.MakeCorrectElectionConstants().Q}
}

func (s *SHA256) update(data string) {
	s.toHash.WriteString(data)
}

func (s *SHA256) digest() *schema.BigInt {
	// Hashing all the data strings
	var hash32 = sha256.Sum256([]byte(s.toHash.String()))

	// Turning byte array into big.Int
	intValueForHash := schema.BigInt{Int: big.Int{}}
	intValueForHash.SetBytes(hash32[:])

	// Taking hash mod q TODO: Should it be q - 1?
	intValueForHash.Mod(&intValueForHash.Int, &s.q.Int)

	return &intValueForHash
}

func HashElems(a ...interface{}) *schema.BigInt {
	h := MakeSHA256()
	h.update("|")

	for _, x := range a {
		var hashMe string

		switch reflect.TypeOf(x) {
		case nilType:
			hashMe = "null"
		case intType:
			// Type cast and take the string representation of the int
			xInt, _ := x.(int)
			hashMe = strconv.Itoa(xInt)
		case stringType:
			// Type cast (strings are already utf8-encoded in Golang)
			hashMe, _ = x.(string)
		case bigIntType:
			// Convert big.Int to hex
			bigInt := x.(schema.BigInt).Int
			hashMe = bigInt.Text(10)
			// Add leading zero if amount of digits is odd // TODO: Might need?
			// hashMe = addLeadingZeroIfNeeded(hex)
		default:
			s := reflect.ValueOf(x)
			var slice = make([]interface{}, s.Len())

			for i := 0; i < s.Len(); i++ {
				slice[i] = s.Index(i).Interface()
			}

			sliceIsEmpty := len(slice) == 0
			if sliceIsEmpty {
				hashMe = "null"
			} else {
				// Else hash the elements and encode to hex
				bigIntRes := HashElems(slice...)
				hashMe = fmt.Sprintf("%X", bigIntRes)
			}
		}
		h.update(hashMe + "|")
	}

	// Digest returns H(strings) mod q
	return h.digest()
}

func addLeadingZeroIfNeeded(hex string) string {
	stringLengthIsOdd := len(hex)%2 != 0
	if stringLengthIsOdd {
		return hex
	}

	return "0" + hex
}
