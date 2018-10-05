package main

import (
	"os"
	"sort"

	"path/filepath"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/csv"
	"github.com/senseyeio/diligent/pretty"
)

type toSortInterfacer func(deps []diligent.Dep) sort.Interface

func getReporter() diligent.Reporter {
	if csvOutput {
		return csv.NewReporter()
	}

	return pretty.NewReporter()
}

func getFiles(args []string) []string {
	path := args[0]
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fatal(66, err.Error())
	}
	return files
}

func run(args []string) {
	files := getFiles(args)

	deps := make([]diligent.Dep, 0)
	warnings := make([]diligent.Warning, 0)
	for _, f := range files {
		deper, err := getDeper(f)
		if err != nil {
			continue
		}
		fileBytes := mustReadFile(f)
		d, w, err := deper.Dependencies(fileBytes)
		if err != nil {
			fatal(67, err.Error())
		}
		deps = append(deps, d...)
		warnings = append(warnings, w...)
	}

	for _, w := range warnings {
		warning(w.Warning())
	}
	if len(deps) == 0 {
		fatal(67, "did not successfully process any dependencies - see warnings above for details")
	}

	deps = diligent.Deps(deps).Dedupe()

	sorter := getSort(sortByLicense)
	sort.Sort(sorter(deps))

	reporter := getReporter()
	if err := reporter.Report(os.Stdout, deps); err != nil {
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
