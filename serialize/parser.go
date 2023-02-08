package serialize

import (
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/util"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseContext() {
	// Open json file and print error if any
	jsonFile, err := os.Open("data/hamilton-general/election_record/context.json")
	util.PrintError(err)

	// Assign variable which the file is to be unmarshalled into
	var cipherTextElectionRecord schema.CiphertextElectionRecord

	// Turn the file into a byte array, and print error
	jsonByte, _ := io.ReadAll(jsonFile)
	jsonErr := json.Unmarshal(jsonByte, &cipherTextElectionRecord)
	util.PrintError(jsonErr)

	// Test print
	fmt.Println(cipherTextElectionRecord.ElgamalPublicKey)

	defer jsonFile.Close()
}
