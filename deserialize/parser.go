package deserialize

import (
	"electionguard-verifier-go/schema"
	"electionguard-verifier-go/utility"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

type ElectionRecord struct {
	// Election data fields
	CiphertextElectionRecord  schema.CiphertextElectionRecord
	Manifest                  schema.Manifest
	ElectionConstants         schema.ElectionConstants
	EncryptedTally            schema.EncryptedTally
	PlaintextTally            schema.PlaintextTally
	CoefficientsValidationSet schema.CoefficientsValidationSet
	SubmittedBallots          []schema.SubmittedBallot
	SpoiledBallots            []schema.SpoiledBallot
	EncryptionDevices         []schema.EncryptionDevice
	Guardians                 []schema.Guardian
}

func MakeElectionRecord() *ElectionRecord {
	return &ElectionRecord{}
}

type Parser struct {
	logger zap.Logger
}

func MakeParser(logger *zap.Logger) *Parser {
	return &Parser{logger: *logger}
}

func (p *Parser) ParseElectionRecord(path string) *ElectionRecord {
	// Creating verifier arguments struct
	verifierArguments := *MakeElectionRecord()

	// Parsing singleton files
	verifierArguments.CiphertextElectionRecord = parseJsonToGoStruct(p.logger, path+"/context.json", schema.CiphertextElectionRecord{})
	verifierArguments.Manifest = parseJsonToGoStruct(p.logger, path+"/manifest.json", schema.Manifest{})
	verifierArguments.EncryptedTally = parseJsonToGoStruct(p.logger, path+"/encrypted_tally.json", schema.EncryptedTally{})
	verifierArguments.ElectionConstants = parseJsonToGoStruct(p.logger, path+"/constants.json", schema.ElectionConstants{})
	verifierArguments.PlaintextTally = parseJsonToGoStruct(p.logger, path+"/tally.json", schema.PlaintextTally{})
	verifierArguments.CoefficientsValidationSet = parseJsonToGoStruct(p.logger, path+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	verifierArguments.EncryptionDevices = parseJsonToSlice(p.logger, path+"/encryption_devices/", schema.EncryptionDevice{})
	verifierArguments.SpoiledBallots = parseJsonToSlice(p.logger, path+"/spoiled_ballots/", schema.SpoiledBallot{})
	verifierArguments.SubmittedBallots = parseJsonToSlice(p.logger, path+"/submitted_ballots/", schema.SubmittedBallot{})
	verifierArguments.Guardians = parseJsonToSlice(p.logger, path+"/guardians/", schema.Guardian{})

	return &verifierArguments
}

func parseJsonToGoStruct[E any](logger zap.Logger, path string, typeOfObject E) E {
	logger.Debug("parsing file: " + path)

	// Open json file and print error if any
	file, fileErr := os.Open(path)
	utility.PanicError(fileErr)

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	utility.PanicError(byteErr)

	// Unmarshal the bytearray into empty instance of variable of type E
	jsonErr := json.Unmarshal(jsonByte, &typeOfObject)
	utility.PanicError(jsonErr)
	if jsonErr != nil {
		fmt.Println(path)
	}

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		utility.PanicError(closeErr)
	}(file)

	return typeOfObject
}

func parseJsonToSlice[E any](logger zap.Logger, path string, typeOfObject E) []E {
	// Getting all files in directory
	files, err := os.ReadDir(path)
	utility.PanicError(err)

	// Creating list and parsing all files in directory
	var l []E
	for _, file := range files {
		var xd E
		toBeAppended := parseJsonToGoStruct(logger, path+file.Name(), &xd)
		l = append(l, *toBeAppended)
	}

	return l
}
