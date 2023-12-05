# ElectionGuard Verifier in Go
*Andreas Skriver Nielsen, Niklas Bille, Markus Valdemar Grønkjær Jensen, and Hans-Christian Kjeldsen*

## Overview
[ElectionGuard](https://github.com/microsoft/electionguard) is an open-source software development kit (SDK) by Microsoft,
that aims to improve the security and transparency of elections. The primary focus is that

- Individual voters can verify that their votes have been accurately recorded.
- Voters and observers can verify that all recorded votes have been accurately counted.

Our independent verifier written in Go, allows voters and observers to confirm the consistency of an election using the supplied election record.

## Installation
As a prerequisite make sure to have installed Go, which can be downloaded [here](https://go.dev/doc/install). Afterwards download the verifier as a ZIP, or clone the repository from source:
```
$ git clone https://github.com/AU-HC/electionguard-verifier-go.git 
```
Then get the dependencies used by the verifier:
```
$ cd electionguard-verifier-go
$ go get
```

## Usage
The verifier is currently a command line utility tool, to verify an election the following command has to be executed.
```
$ go run main.go -p="path/to/election-record/"
```
It's important to note that the `-p` flag must be set, as it specifies the election record path. The election record
must follow the specification of ElectionGuard version `1.91.18`.

The verifier also has alternate options which can be set, using the following flags:
- `-o` of type `string`: Specifies if the verifier, should output a JSON file with additional verification information to the specified path.
- `-c` of type `bool` : Specifies if the verifier should use multiple cores to verify the election record.
- `-b` of type `int`: Specifies the amount of samples for a benchmarking run. Setting this flag with a value other than 0, will not verify the specified election.
- `-v` of type `int`: Specifies the logging level for the verifier, the options are:
    - *0* : Will log nothing (default)
    - *1* : Logging of information
    - *2* : Logging of debug

The project provides one sample election in the `/data/nov_2023` directory, by courtesy of [ElectionGuard](https://www.electionguard.vote/elections/College_Park_Maryland_2023/). 
To verify the College Park General Election with logging level set to `information` and output file `output.json` execute one of the following command blocks.
```
$ go run main.go -p="data/nov_2023/election_record/" -v=1 -o="output.json" 
```
or (Windows)
```
$ go build main.go
$ electionguard-verifier-go.exe -p="data/nov_2023/election_record/" -v=1 -o="output.json" 
```
or (Mac/Linux)
```
$ go build main.go
$ ./electionguard-verifier-go -p="data/nov_2023/election_record/" -v=1 -o="output.json" 
```

## Remarks
- The verifier is currently verifying specification `1.91.18`, a hybrid version that combines specification `2.0` and `1.53`. This specification aligns with the outlined criteria detailed in the [requirements document](https://www.electionguard.vote/images/MITRE-EG-CP-requirements.pdf) by MITRE.
- To verify an `1.0` or `1.1` specification election please use a previous version of the verifier, which can be found [here](https://github.com/AU-HC/electionguard-verifier-go/tree/main/version/1.1).