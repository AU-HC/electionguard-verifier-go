# ElectionGuard Verifier in Go
*Andreas Skriver Nielsen, Niklas Bille Olesen, and Hans-Christian Kjeldsen*

## Overview
[ElectionGuard](https://github.com/microsoft/electionguard) is an open-source software development kit (SDK) by Microsoft,
that aims to improve the security and transparency of elections. The primary focus is that

- Individual voters can verify that their votes have been accurately recorded.
- Voters and observers can verify that all recorded votes have been accurately counted.

Our independent verifier written in Go, allows voters and observers to confirm the consistency of an election using the supplied election artifacts
along with the published election results.

## Installation
As a prerequisite make sure to have installed
- Go, which can be downloaded [here](https://go.dev/doc/install)

Download the verifier as a ZIP, or clone the repository from source:
```
$ git clone https://github.com/AU-HC/electionguard-verifier-go.git 
```

## Usage
The verifier is currently a command line utility tool, to verify an election the following command has to be executed.
```
$ go run main.go --p="path/to/election-record/"
```
It's important to note that the `-p` flag must be set, as it specifies the election record path. The election record
must follow the specification of ElectionGuard version 1.0 or 1.1

The verifier also has alternate options which can be set, using the following flags:
- `-o` of type `string`: Specifies if the verifier, should output a JSON file with additional verification information to the specified path.
- `-c` of type `bool` : Specifies if the verifier should use multiple cores to verify the election record.
- `-b` of type `int`: Specifies the amount of samples for a benchmarking run. Setting this flag with a value other than 0, will not verify the specified election.
- `-v` of type `int`: Specifies the logging level for the verifier, the options are:
    - *0* : Will log nothing (default)
    - *1* : Logging of information
    - *2* : Logging of debug

The project provides some sample data in `/data`, courtesy of [Microsoft](https://github.com/microsoft/electionguard/releases/tag/v1.1) and [egvote.us](https://www.egvote.us/cc/id/22). 
To verify `/data/full/` with logging level set to `information` and output file `output.json` execute the following command
```
$ go run main.go --p="data/full/election_record/" --v=1 --o="output.json" 
```
or (Windows)
```
$ go build main.go
$ electionguard-verifier-go.exe --p="data/full/election_record/" --v=1 --o="output.json" 
```
or (Mac/Linux)
```
$ go build main.go
$ ./electionguard-verifier-go --p="data/full/election_record/" --v=1 --o="output.json" 
```

## Additional information
### Remarks on the verifier
* The verifier is currently not verifying step `6.A` as the ElectionGuard specification is not detailed enough.
* The verifier is currently not verifying step `11.E` and `11.F` because the information needed to verify is not located in spoiled ballots.
* We have not yet been able to verify that step `17`, `18`, and `19` is being verified correctly as no election record contains the artifacts to verify the steps. 

### Sample data
* All the sample data can be found on the release.
* Some of the sample data is already included in the project folder `/data`.

### Backlog
- [ ] Comment code, and general cleanup
- [ ] Upload report to GitHub, and create section in `README.md`
