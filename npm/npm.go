package npm

import (
	"encoding/json"
	"fmt"
	"github.com/senseyeio/diligent"
	"io/ioutil"
	"net/http"
	"strings"
)

type packageJson struct {
	Deps map[string]string `json:"dependencies"`
}

type npmPackage struct {
	License string `json:"license"`
}

type npmDeper struct{}

func New() diligent.Deper {
	return &npmDeper{}
}

func (n *npmDeper) Name() string {
	return "npm"
}

func (n *npmDeper) Dependencies(file []byte, options map[string]interface{}) ([]diligent.Dep, error) {
	var pkg packageJson
	err := json.Unmarshal(file, &pkg)
	if err != nil {
		return nil, err
	}

	deps := make([]diligent.Dep, 0, len(pkg.Deps))
	for pkg, version := range pkg.Deps {
		l, err := getNPMLicense(pkg, version)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to get license for %s: %s", pkg, err.Error()))
		} else {
			deps = append(deps, l)
		}
	}
	return deps, nil
}
func (n *npmDeper) IsCompatible(filename string, fileContents []byte) bool {
	return strings.Index(filename, "package.json") != -1
}

func getNPMLicense(pkgName, version string) (diligent.Dep, error) {
	resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/%s/%s", pkgName, version))
	if err != nil {
		return diligent.Dep{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diligent.Dep{}, err
	}

	var packageInfo npmPackage
	err = json.Unmarshal(body, &packageInfo)
	if err != nil {
		return diligent.Dep{}, err
	}

	return diligent.Dep{
		Name:    pkgName,
		License: diligent.GetLicenseFromIdentifier(packageInfo.License),
	}, nil
}
