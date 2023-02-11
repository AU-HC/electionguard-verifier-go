package serialize

import (
	"electionguard-verifier-go/utility"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

type Parser struct {
	logger zap.Logger
}

func MakeParser(logger *zap.Logger) *Parser {
	return &Parser{logger: *logger}
}

func ParseFromJsonToSingleObject[E any](path string, typeOfObject E) E {
	// Open json file and print error if any
	file, fileErr := os.Open(path)
	utility.PrintError(fileErr)

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	utility.PrintError(byteErr)

	// Unmarshal the bytearray into empty instance of variable of type E
	// and print any error
	jsonErr := json.Unmarshal(jsonByte, &typeOfObject)
	utility.PrintError(jsonErr)
	if jsonErr != nil {
		fmt.Println(path)
	}

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		utility.PrintError(closeErr)
	}(file)

	return typeOfObject
}

func ParseFromJsonToSlice[E any](path string, typeOfObject E) []E {
	// Getting all files in directory
	files, err := os.ReadDir(path)
	utility.PrintError(err)

	// Creating list and parsing all files in directory
	var l []E
	for _, file := range files {
		toBeAppended := ParseFromJsonToSingleObject(path+file.Name(), &typeOfObject)
		l = append(l, *toBeAppended)
	}

	return l
}
