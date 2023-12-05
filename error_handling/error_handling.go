package error_handling

// PanicError is used to halt the given goroutine if there is an error present
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
