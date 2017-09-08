package npm

import (
	"encoding/json"
	"fmt"
	"github.com/senseyeio/diligent"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type packageJson struct {
	Deps    map[string]string `json:"dependencies"`
	DevDeps map[string]string `json:"devDependencies"`
}

type npmPackage struct {
	License string `json:"license"`
}

type npmDeper struct {
	config Config
}

type Config struct {
	DevDependencies bool
}

func New() diligent.Deper {
	return NewWithOptions(Config{})
}

func NewWithOptions(c Config) diligent.Deper {
	return &npmDeper{c}
}

func (n *npmDeper) Name() string {
	return "npm"
}

func mergeMaps(to map[string]string, from map[string]string) {
	for pkg, version := range from {
		to[pkg] = version
	}
}

func (n *npmDeper) Dependencies(file []byte, options map[string]interface{}) ([]diligent.Dep, error) {
	var pkg packageJson
	err := json.Unmarshal(file, &pkg)
	if err != nil {
		return nil, err
	}

	licensesToGet := map[string]string{}
	mergeMaps(licensesToGet, pkg.Deps)
	if n.config.DevDependencies {
		mergeMaps(licensesToGet, pkg.DevDeps)
	}

	deps := make([]diligent.Dep, 0, len(licensesToGet))
	for pkg, version := range licensesToGet {
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
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", url.QueryEscape(pkgName), url.QueryEscape(version))
	resp, err := http.Get(url)
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
