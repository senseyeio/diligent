package govendor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/go"
)

type pkg struct {
	Path     string `json:"path"`
	Revision string `json:"revision"`
}

type vendor struct {
	Packages []pkg `json:"package"`
}

type govendor struct{}

func New() diligent.Deper {
	return &govendor{}
}

func (g *govendor) Name() string {
	return "govendor"
}

func (g *govendor) Dependencies(file []byte) ([]diligent.Dep, error) {
	var vendorFile vendor
	err := json.Unmarshal(file, &vendorFile)
	if err != nil {
		log.Fatal(err)
	}

	deps := make([]diligent.Dep, 0, len(vendorFile.Packages))
	for _, pkg := range vendorFile.Packages {
		pkgPath := pkg.Path
		l, err := _go.GetLicense(pkgPath)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to get license for %s: %s", pkgPath, err.Error()))
		} else {
			deps = append(deps, diligent.Dep{
				Name:    pkgPath,
				License: l,
			})
		}
	}
	return deps, nil
}
func (g *govendor) IsCompatible(filename string, fileContents []byte) bool {
	return strings.Index(filename, "vendor.json") != -1
}
