package govendor

import (
	"encoding/json"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/warning"
)

type pkg struct {
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

type vendor struct {
	Packages []pkg `json:"package"`
}

type govendor struct {
	lg GoLicenseGetter
}

type GoLicenseGetter interface {
	GetLicense(packagePath string) (diligent.License, error)
}

// New returns a Deper capable of handling govendor manifest files
func New(lg GoLicenseGetter) diligent.Deper {
	return &govendor{lg}
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
		l, err := g.lg.GetLicense(pkgPath)
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
func (g *govendor) IsCompatible(filename string) bool {
	return filename == "vendor.json"
}
