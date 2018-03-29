package govendor

import (
	"encoding/json"
	"strings"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/go"
	"github.com/senseyeio/diligent/warning"
)

type pkg struct {
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

type vendor struct {
	Packages []pkg `json:"package"`
}

type govendor struct{}

// New returns a Deper capable of handling govendor manifest files
func New() diligent.Deper {
	return &govendor{}
}

// Name returns "govendor"
func (g *govendor) Name() string {
	return "govendor"
}

// Dependencies returns the licenses of the go packages defined within the govendor manifest
func (g *govendor) Dependencies(file []byte) ([]diligent.Dep, []diligent.Warning, error) {
	var vendorFile vendor
	err := json.Unmarshal(file, &vendorFile)
	if err != nil {
		return nil, nil, err
	}

	deps := make([]diligent.Dep, 0, len(vendorFile.Packages))
	warns := make([]diligent.Warning, 0, len(vendorFile.Packages))
	for _, pkg := range vendorFile.Packages {
		pkgPath := pkg.Path
		l, err := _go.GetLicense(pkgPath)
		if err != nil {
			warns = append(warns, warning.New(pkgPath, err.Error()))
		} else {
			deps = append(deps, diligent.Dep{
				Name:    pkgPath,
				License: l,
			})
		}
	}
	return deps, warns, nil
}

// IsCompatible returns true if the filename is vendor.json
func (g *govendor) IsCompatible(filename string, fileContents []byte) bool {
	return strings.Index(filename, "vendor.json") != -1
}
