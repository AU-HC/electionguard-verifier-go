package schema

import (
	"fmt"
	"math/big"
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
	_, ok := z.SetString(s, 16) // TODO: Check base
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
