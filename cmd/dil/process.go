package main

import (
	"github.com/senseyeio/diligent"
	"io/ioutil"
	"log"
)

func runDep(deper diligent.Deper, reper diligent.Reporter, filePath string) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	deps, err := deper.Dependencies(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = reper.Report(deps); err != nil {
		log.Fatal(err.Error())
	}

	if err = validateDependencies(deps); err != nil {
		log.Fatal(err.Error())
	}
}
