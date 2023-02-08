package serialize

import (
	"electionguard-verifier-go/util"
	"encoding/json"
	"io"
	"os"
)

func ParseFromJson[E any](path string, t E) E {
	// Open json file and print error if any
	file, fileErr := os.Open(path)
	util.PrintError(fileErr)

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	util.PrintError(byteErr)

	// Unmarshal the bytearray into empty instance of variable of type E
	// and print any error
	jsonErr := json.Unmarshal(jsonByte, &t)
	util.PrintError(jsonErr)

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		util.PrintError(closeErr)
	}(file)

	return t
}
