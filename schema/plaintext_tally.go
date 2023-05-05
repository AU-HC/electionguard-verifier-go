package schema

type PlaintextTally struct {
	ObjectId string                  `json:"object_id"`
	Contests map[string]ContestTally `json:"contests"`
}

type ContestTally struct {
	ObjectId   string                    `json:"object_id"`
	Selections map[string]SelectionTally `json:"selections"`
}

type SelectionTally struct {
	ObjectId string           `json:"object_id"`
	Tally    int              `json:"tally"`
	Value    BigInt           `json:"value"`
	Message  Ciphertext       `json:"message"`
	Shares   []SelectionShare `json:"shares"`
}

type SelectionShare struct {
	ObjectId       string                   `json:"object_id"`
	GuardianId     string                   `json:"guardian_id"`
	Share          BigInt                   `json:"share"`
	Proof          CpProof                  `json:"proof"`
	RecoveredParts map[string]RecoveredPart `json:"recovered_parts"`
}

type CpProof struct {
	Pad       BigInt `json:"pad"`
	Data      BigInt `json:"data"`
	Challenge BigInt `json:"challenge"`
	Response  BigInt `json:"response"`
	Usage     string `json:"usage"`
}

type RecoveredPart struct {
	ObjectId                  string  `json:"object_id"`
	GuardianIdentifier        string  `json:"guardian_id"`
	MissingGuardianIdentifier string  `json:"missing_guardian_id"`
	Share                     BigInt  `json:"share"`
	RecoveryPublicKey         BigInt  `json:"recovery_key"`
	Proof                     CpProof `json:"proof"`
}

func (c *CpProof) IsNotEmpty() bool {
	zero := MakeBigIntFromInt(0)

	if c.Pad.Compare(zero) {
		return false
	}

	if c.Data.Compare(zero) {
		return false
	}

	if c.Challenge.Compare(zero) {
		return false
	}

	if c.Response.Compare(zero) {
		return false
	}

	if c.Usage == "" {
		return false
	}

	return true
}

func (s *SelectionShare) IsNotEmpty() bool {
	if s.ObjectId == "" {
		return false
	}

	if s.GuardianId == "" {
		return false
	}

	if s.Share.Compare(MakeBigIntFromInt(0)) {
		return false
	}

	if s.ObjectId == "" {
		return false
	}

	return true
}

func (r *RecoveredPart) IsNotEmpty() bool {
	if r.ObjectId == "" {
		return false
	}

	if r.GuardianIdentifier == "" {
		return false
	}

	if r.MissingGuardianIdentifier == "" {
		return false
	}

	if r.Share.Compare(MakeBigIntFromInt(0)) {
		return false
	}

	if r.RecoveryPublicKey.Compare(MakeBigIntFromInt(0)) {
		return false
	}

	return true
}
