package main

import (
	"fmt"
	"github.com/senseyeio/diligent"
	"io/ioutil"
	"log"
)

func runDep(deper diligent.Deper, filePath string) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(deper.Name())

	deps, err := deper.Dependencies(fileBytes, map[string]interface{}{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, dep := range deps {
		fmt.Println(fmt.Sprintf("%s -> %s", dep.Name, dep.License.Name))
	}
}
