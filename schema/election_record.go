package schema

type ElectionRecord struct {
	// Election data fields
	CiphertextElectionRecord  Context
	Manifest                  Manifest
	ElectionConstants         ElectionConstants
	EncryptedTally            EncryptedTally
	PlaintextTally            PlaintextTally
	CoefficientsValidationSet CoefficientsValidationSet
	SubmittedBallots          []SubmittedBallot
	MockBallots               []MockBallot
	SpoiledBallots            []SpoiledBallot
	EncryptionDevices         []EncryptionDevice
	Guardians                 []Guardian
}
