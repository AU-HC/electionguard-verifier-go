package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"electionguard-verifier-go/schema"
	"encoding/binary"
	"fmt"
	"reflect"
)

func HMAC(key schema.BigInt, domainSeparator byte, a ...interface{}) *schema.BigInt {
	mac := hmac.New(sha256.New, key.Bytes()) // should test this actually appends as expected

	// Add the domain separator first
	mac.Write([]byte{domainSeparator})

	// Then append the message (i.e. what is to be hashed)
	for _, x := range a {
		var toBeHashed []byte

		switch reflect.TypeOf(x) { // TODO: add files
		case intType:
			// Type cast and create byte array which the number is to be stored in
			xInt, _ := x.(int)

			// We know that all small integers in ElectionGuard is smaller than 2^31, therefore we can typecast to uint32
			toBeHashed = make([]byte, 4)
			binary.BigEndian.PutUint32(toBeHashed, uint32(xInt))

		case stringType:
			// Type cast (strings are already utf8-encoded in Golang)
			xString, _ := x.(string)

			// Pad the string as a byte with four empty bytes
			pad := make([]byte, 4)
			toBeHashed = append(pad, []byte(xString)...)

		case bigIntType:
			// Convert big.Int to hex
			bigInt := x.(schema.BigInt).Int
			fmt.Println("Length of bigint: (should be 32 or 512)", len(bigInt.Bytes()))
			toBeHashed = bigInt.Bytes()
		}

		mac.Write(toBeHashed)
	}

	// TODO: Should you also take mod q of the result
	hash := schema.MakeBigIntFromByteArray(mac.Sum(nil))

	return hash
}
