package utility

import "fmt"

// PrintError is used to handle go errors throughout the project
func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// PanicError is used to halt the given goroutine if there is an error present
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
