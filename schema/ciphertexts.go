package schema

type ExtendedData struct {
	Generator                 BigInt `json:"pad"`
	EncryptedMessage          BigInt `json:"data"`
	MessageAuthenticationCode BigInt `json:"mac"`
}
