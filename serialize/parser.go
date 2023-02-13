package serialize

import (
	"electionguard-verifier-go/core"
	"electionguard-verifier-go/schema"
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

func (p *Parser) ConvertJsonDataToGoStruct(path string) core.VerifierArguments {
	// Creating verifier arguments struct
	verifierArguments := *core.MakeVerifierArguments()

	// Parsing singleton files
	verifierArguments.CiphertextElectionRecord = parseJsonStruct(p.logger, path+"/context.json", schema.CiphertextElectionRecord{})
	verifierArguments.Manifest = parseJsonStruct(p.logger, path+"/manifest.json", schema.Manifest{})
	verifierArguments.EncryptedTally = parseJsonStruct(p.logger, path+"/encrypted_tally.json", schema.EncryptedTally{})
	verifierArguments.ElectionConstants = parseJsonStruct(p.logger, path+"/constants.json", schema.ElectionConstants{})
	verifierArguments.PlaintextTally = parseJsonStruct(p.logger, path+"/tally.json", schema.PlaintextTally{})
	verifierArguments.CoefficientsValidationSet = parseJsonStruct(p.logger, path+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	verifierArguments.EncryptionDevices = parseJsonToSlice(p.logger, path+"/encryption_devices/", schema.EncryptionDevice{})
	verifierArguments.SpoiledBallots = parseJsonToSlice(p.logger, path+"/spoiled_ballots/", schema.SpoiledBallot{})
	verifierArguments.SubmittedBallots = parseJsonToSlice(p.logger, path+"/submitted_ballots/", schema.SubmittedBallot{})
	verifierArguments.Guardians = parseJsonToSlice(p.logger, path+"/guardians/", schema.Guardian{})

	return verifierArguments
}

func parseJsonStruct[E any](logger zap.Logger, path string, typeOfObject E) E {
	logger.Debug("parsing file: " + path)

	// Open json file and print error if any
	file, fileErr := os.Open(path)
	utility.PrintError(fileErr)

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	utility.PrintError(byteErr)

	// Unmarshal the bytearray into empty instance of variable of type E
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

func parseJsonToSlice[E any](logger zap.Logger, path string, typeOfObject E) []E {
	// Getting all files in directory
	files, err := os.ReadDir(path)
	utility.PrintError(err)

	// Creating list and parsing all files in directory
	var l []E
	for _, file := range files {
		toBeAppended := parseJsonStruct(logger, path+file.Name(), &typeOfObject)
		l = append(l, *toBeAppended)
	}

	return l
}
