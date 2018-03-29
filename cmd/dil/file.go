package main

import "io/ioutil"

func mustReadFile(path string) []byte {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fatal(66, err.Error())
	}
	return fileBytes
}
