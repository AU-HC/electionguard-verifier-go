package serialize

import (
	"electionguard-verifier-go/util"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseFromJsonToSingleObject[E any](path string, t E) E {
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

	if jsonErr != nil {
		fmt.Print(path)
	}

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		util.PrintError(closeErr)
	}(file)

	return t
}

func ParseFromJsonToSlice[E any](path string, t E) []E {
	// Getting all files in directory
	files, err := os.ReadDir(path)
	util.PrintError(err)

	// Creating list and parsing all files in directory
	var l []E
	for _, file := range files {
		toBeAppended := ParseFromJsonToSingleObject(path+file.Name(), &t)
		l = append(l, *toBeAppended)
	}

	return l
}
