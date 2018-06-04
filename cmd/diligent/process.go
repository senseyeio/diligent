package main

import (
	"github.com/senseyeio/diligent"
	"os"
)

func runDep(deper diligent.Deper, reporter diligent.Reporter, filePath string) {
	fileBytes := mustReadFile(filePath)
	deps, warnings, err := deper.Dependencies(fileBytes)
	if err != nil {
		fatal(67, err.Error())
	}

	for _, w := range warnings {
		warning(w.Warning())
	}
	if len(deps) == 0 {
		fatal(67, "did not successfully process any dependencies - see warnings above for details")
	}

	if err = reporter.Report(deps); err != nil {
		fatal(65, err.Error())
	}

	if err = validateDependencies(deps); err != nil {
		fatal(68, err.Error())
	}

	if len(warnings) > 0 {
		os.Exit(64)
	}
}