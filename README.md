# ElectionGuard Verifier in Go
...

## Installation
...

## Usage
...

## TODO
- [x] Verify step 14
- [x] Double check `schema/manifest.go`
- [x] Refactor `schema/manifest.go`
- [x] Verify step 18
- [x] Verify step 19
- [x] Verify step 7
- [x] General refactor of `core/verifier.go`
- [x] Do proper error handling when parsing JSON data
- [x] Confirm step 17 works according to spec sheet
- [x] Add parallelization
- [x] Make it possible to get output as a file with information (such as amount of checked invariants for each step, etc.)
- [x] Add parallelization for step 4, 5, 9 (Split slice into n slices)
- [x] Refactor validation helper in `core/validate_utility.go`
- [x] Verify step 16C to 16E
- [ ] Refactor WaitGroup in `core/verifier.go` to not have `wg.Add(1)` in goroutines
- [ ] Check 16.B
- [ ] Verify step 11C to 11F
- [ ] Verify step 6A
- [ ] Finish `README.md`
- [ ] Is `schema/manifest.go` supposed to be bricked? Does it even matter?
