package schema

type ElectionConstants struct {
	LargePrime BigInt `json:"large_prime"`
	SmallPrime BigInt `json:"small_prime"`
	Cofactor   BigInt `json:"cofactor"`
	Generator  BigInt `json:"generator"`
}
