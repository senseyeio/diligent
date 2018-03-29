package main

import (
	"github.com/senseyeio/diligent"
	"fmt"
	"os"
)

func runDep(deper diligent.Deper, reper diligent.Reporter, filePath string) {
	fileBytes := mustReadFile(filePath)
	deps, warnings, err := deper.Dependencies(fileBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(67)
	}

	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, w.Warning())
	}

	if err = reper.Report(deps); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(65)
	}

	if err = validateDependencies(deps); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(68)
	}

	if len(warnings) > 0 {
		os.Exit(64)
	}
}
