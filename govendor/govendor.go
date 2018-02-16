package govendor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ryanuber/go-license"
	"github.com/senseyeio/diligent"
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
		l, err := getGoLicense(pkgPath)
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

func getGoLicense(goPackagePath string) (diligent.License, error) {
	cmd := exec.Command("go", "get", goPackagePath)
	cmd.Run()

	l, err := license.NewFromDir(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), goPackagePath))
	if err != nil {
		components := strings.Split(goPackagePath, "/")
		if len(components) > 3 {
			return getGoLicense(strings.Join(components[:len(components)-1], "/"))
		}
		return diligent.License{}, err
	}

	return diligent.GetLicenseFromIdentifier(l.Type)
}
