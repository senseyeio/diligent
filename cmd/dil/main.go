package main

import (
	"fmt"
	"github.com/senseyeio/diligent"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	fileName := args[0]

	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	deper, err := getDeper(fileName, fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(deper.Name())

	deps, err := deper.Dependencies(fileBytes, map[string]interface{}{})
	if err != nil {
		log.Fatal(err.Error())
	}

	printDeps(deps)
}

func printDeps(deps []diligent.Dep) {
	for _, dep := range deps {
		fmt.Println(fmt.Sprintf("%s -> %s", dep.Name, dep.License.Name))
	}
}
