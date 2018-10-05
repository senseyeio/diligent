package main

import (
	"os"
	"sort"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/csv"
	"github.com/senseyeio/diligent/pretty"
)

type toSortInterfacer func(deps []diligent.Dep) sort.Interface

func getReporter() diligent.Reporter {
	if csvFilePath != "" {
		return csv.NewReporter(csvFilePath)
	}

	return pretty.NewReporter()
}

func run(args []string) {
	filePath := args[0]
	fileBytes := mustReadFile(filePath)
	deper, err := getDeper(filePath, fileBytes)
	if err != nil {
		fatal(69, err.Error())
	}

	runDep(deper, getSort(sortByLicense), getReporter(), filePath)
}

func runDep(deper diligent.Deper, sorter toSortInterfacer, reporter diligent.Reporter, filePath string) {
	fileBytes := mustReadFile(filePath)
	deps, warnings, err := deper.Dependencies(fileBytes)
	if err != nil {
		fatal(67, err.Error())
	}

	sort.Sort(sorter(deps))

	for _, w := range warnings {
		warning(w.Warning())
	}
	if len(deps) == 0 {
		fatal(67, "did not successfully process any dependencies - see warnings above for details")
	}

	if err = reporter.Report(deps); err != nil {
		fatal(65, err.Error())
	}

	if errs := validateDependencies(deps); len(errs) > 0 {
		if len(errs) == 1 {
			fatal(68, errs[0].Error())
		}
		for _, e := range errs {
			warning(e.Error())
		}
		fatal(68, "multiple dependencies are not compliant with your whitelist")
	}

	if len(warnings) > 0 {
		os.Exit(64)
	}
}

func toLicenseSorter(deps []diligent.Dep) sort.Interface {
	return diligent.DepsByLicense(deps)
}

func toNameSorter(deps []diligent.Dep) sort.Interface {
	return diligent.DepsByName(deps)
}

func getSort(useLicenseSorting bool) toSortInterfacer {
	if useLicenseSorting {
		return toLicenseSorter
	}

	return toNameSorter
}
