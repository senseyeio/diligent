package dep

import (
	"fmt"
	"github.com/ryanuber/go-license"
	"github.com/senseyeio/diligent"
	"github.com/pelletier/go-toml"
	"log"
	"os"
	"os/exec"
	"strings"
)

type lockedProject struct {
	Name     string   `toml:"name"`
}

type lock struct {
	Projects  []lockedProject `toml:"projects"`
}

type dep struct{}

func New() diligent.Deper {
	return &dep{}
}

func (d *dep) Name() string {
	return "dep"
}

func (d *dep) Dependencies(file []byte) ([]diligent.Dep, error) {
	var l lock
	err := toml.Unmarshal(file, &l)
	if err != nil {
		log.Fatal(err)
	}

	deps := make([]diligent.Dep, 0, len(l.Projects))
	for _, pkg := range l.Projects {
		l, err := getGoLicense(pkg.Name)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to get license for %s: %s", pkg.Name, err.Error()))
		} else {
			deps = append(deps, diligent.Dep{
				Name:    pkg.Name,
				License: l,
			})
		}
	}
	return deps, nil
}
func (d *dep) IsCompatible(filename string, fileContents []byte) bool {
	return strings.Index(filename, "Gopkg.lock") != -1
}

func getGoLicense(goPackagePath string) (diligent.License, error) {
	cmd := exec.Command("go", "get", goPackagePath)
	cmd.Run()

	l, err := license.NewFromDir(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), goPackagePath))
	if err != nil {
		return diligent.License{}, err
	}

	return diligent.GetLicenseFromIdentifier(l.Type), nil
}
