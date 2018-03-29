package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

func mustReadFile(path string) []byte {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(66)
	}
	return fileBytes
}