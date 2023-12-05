package crypto

import (
	"bytes"
	"crypto/sha256"
	"electionguard-verifier-go/schema"
	"reflect"
	"strconv"
	"strings"
)

var separatorString = "|"

var nilType = reflect.TypeOf(nil)
var stringType = reflect.TypeOf("")
var intType = reflect.TypeOf(1)
var bigIntType = reflect.TypeOf(schema.BigInt{})
var fileType = reflect.TypeOf(([]byte)(nil))
var bigIntPointerType = reflect.TypeOf(schema.MakeBigIntFromString("0", 10))

func Hash(q *schema.BigInt, a ...interface{}) *schema.BigInt {
	var buffer bytes.Buffer

	// Then append the message (i.e. what is to be hashed)
	for _, x := range a {
		// Append the separatorString
		buffer.Write([]byte(separatorString))

		var toBeHashed []byte
		switch reflect.TypeOf(x) {
		case intType:
			// Type cast and create byte array which the number is to be stored in
			xInt, _ := x.(int)
			toBeHashed = []byte(strconv.Itoa(xInt))

		case stringType:
			// Type cast (strings are already utf8-encoded in Golang)
			xString, _ := x.(string)
			toBeHashed = []byte(xString)

		case bigIntType:
			bigInt := x.(schema.BigInt)
			toBeHashed = hashBigInt(&bigInt)

		case bigIntPointerType:
			bigIntPointer := x.(*schema.BigInt)
			toBeHashed = hashBigInt(bigIntPointer)

		default:
			panic("unknown type for hash")
		}

		buffer.Write(toBeHashed)
	}
	// Append the separatorString
	buffer.Write([]byte(separatorString))

	// Hash the input and take mod q
	hashByteArray := sha256.Sum256(buffer.Bytes())
	hash := schema.MakeBigIntFromByteArray(hashByteArray[:])
	hash.Mod(&hash.Int, &q.Int)
	return hash
}

func hashBigInt(bigInt *schema.BigInt) []byte {
	bigIntString := bigInt.Text(16)

	if len(bigIntString)%2 == 1 {
		bigIntString = "0" + bigIntString
	}
	bigIntString = strings.ToUpper(bigIntString)
	return []byte(bigIntString)
}
