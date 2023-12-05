package deserialize

import (
	"electionguard-verifier-go/error_handling"
	"electionguard-verifier-go/schema"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
)

type ElectionRecord struct {
	// Election data fields
	CiphertextElectionRecord  schema.Context
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
	logger   *zap.Logger
	errorMsg *strings.Builder
}

func MakeParser(logger *zap.Logger) *Parser {
	return &Parser{logger: logger, errorMsg: &strings.Builder{}}
}

func (p *Parser) ParseElectionRecord(path string) (*ElectionRecord, string) {
	// Creating election record struct
	electionRecord := *MakeElectionRecord()

	// Parsing singleton files
	electionRecord.CiphertextElectionRecord = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/context.json", schema.Context{})
	electionRecord.Manifest = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/manifest.json", schema.Manifest{})
	electionRecord.EncryptedTally = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/encrypted_tally.json", schema.EncryptedTally{})
	electionRecord.ElectionConstants = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/constants.json", schema.ElectionConstants{})
	electionRecord.PlaintextTally = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/tally.json", schema.PlaintextTally{})
	electionRecord.CoefficientsValidationSet = parseJsonToGoStruct(p.logger, p.errorMsg, path+"/coefficients.json", schema.CoefficientsValidationSet{})

	// Directory of file(s)
	electionRecord.EncryptionDevices = parseJsonToSlice(p.logger, p.errorMsg, path+"/encryption_devices/", schema.EncryptionDevice{})
	electionRecord.SpoiledBallots = parseJsonToSlice(p.logger, p.errorMsg, path+"/spoiled_ballots/", schema.SpoiledBallot{})
	electionRecord.SubmittedBallots = parseJsonToSlice(p.logger, p.errorMsg, path+"/submitted_ballots/", schema.SubmittedBallot{})
	electionRecord.Guardians = parseJsonToSlice(p.logger, p.errorMsg, path+"/guardians/", schema.Guardian{})

	return &electionRecord, p.errorMsg.String()
}

func parseJsonToGoStruct[E any](logger *zap.Logger, errorMsg *strings.Builder, path string, typeOfObject E) E {
	logger.Debug("parsing file: " + path)

	// Open json file and print error if any
	file, fileErr := os.Open(path)
	if fileErr != nil {
		errorMsg.WriteString("Could not find file at " + path)
		return typeOfObject
	}

	// Turn the file into a byte array, and print error if any
	jsonByte, byteErr := io.ReadAll(file)
	if byteErr != nil {
		errorMsg.WriteString("Could not read from file at " + path)
		return typeOfObject
	}

	// Unmarshal the bytearray into empty instance of variable of type E
	jsonErr := json.Unmarshal(jsonByte, &typeOfObject)
	if jsonErr != nil {
		errorMsg.WriteString("Could not unmarshall file at " + path + ". Got error: " + jsonErr.Error())
		return typeOfObject
	}

	// Defer close on file, and handling any error
	defer func(file *os.File) {
		closeErr := file.Close()
		error_handling.PanicError(closeErr)
	}(file)

	return typeOfObject
}

func parseJsonToSlice[E any](logger *zap.Logger, errorMsg *strings.Builder, path string, typeOfObject E) []E {
	// Getting all files in directory
	files, err := os.ReadDir(path)
	if err != nil {
		errorMsg.WriteString("Could not find folder at " + path)
		return []E{}
	}

	// Creating list and parsing all files in directory
	var l []E
	for _, file := range files {
		var xd E
		toBeAppended := parseJsonToGoStruct(logger, errorMsg, path+file.Name(), &xd)
		l = append(l, *toBeAppended)
	}

	return l
}
