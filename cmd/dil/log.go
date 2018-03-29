package main

import (
	"fmt"
	"os"
)

func fatal(code int, message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(code)
}

func warning(message string) {
	fmt.Fprintln(os.Stderr, message)
}
