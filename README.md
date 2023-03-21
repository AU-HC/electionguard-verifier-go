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
- `-o` of type `string`: Which specifies if the verifier, should output a JSON file with additional verification information to the specified path.
- `-v` of type `int`: Which specifies the logging level for the verifier, the options are:
    - *0* : Will log nothing (default)
    - *1* : Logging of information
    - *2* : Logging of debug

The project provides some sample data in `/data`, which is taken from [Microsoft](https://github.com/microsoft/electionguard/releases/tag/v1.1) and [egvote.us](https://www.egvote.us/cc/id/22). 
To verify `/data/idaho_pilot_2022/` with logging level set to `information` and output file `output.json` execute the following command
```
$ go run main.go --p="data/idaho_pilot_2022/election-record/" --v=1 --o="output.json" 
```

## Remarks
### Note
The verifier is currently not verifying step `6.A` as the ElectionGuard specification is not detailed enough.

### Backlog
- [x] Finish `README.md`
- [x] Check `10.A`, `14.A`
- [ ] Verify step `6A` (Correct confirmation codes)
- [ ] Upload report to GitHub, and create section in `README.md`