package schema

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type BigInt struct {
	big.Int
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	// If string is null, return nil
	if string(p) == "null" {
		return nil
	}

	s := convertByteArrayToStringAndTrim(p)

	// Declare new bigInt and return error if SetString fails
	z := big.Int{}
	_, ok := z.SetString(s, 16)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z

	// Return no error
	return nil
}

func convertByteArrayToStringAndTrim(p []byte) string {
	s := string(p)
	s = strings.TrimRight(s, "\"")
	s = strings.TrimLeft(s, "\"")

	return s
}

func MakeBigIntFromString(data string, base int) *BigInt {
	var result BigInt
	result.SetString(data, base)

	return &result
}

func IntToBigInt(data int) *BigInt {
	var result BigInt
	stringOfData := strconv.Itoa(data)
	result.SetString(stringOfData, 10)

	return &result
}

func MakeBigIntFromByteArray(bytes []byte) *BigInt {
	var result BigInt
	result.SetBytes(bytes)

	return &result
}

// Compare Method to make true/false comparisons easier
func (b *BigInt) Compare(a *BigInt) bool {
	isEqual := b.Cmp(&a.Int)
	return isEqual == 0
}
