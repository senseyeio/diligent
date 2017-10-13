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

func (n *npmDeper) Dependencies(file []byte) ([]diligent.Dep, error) {
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

func getNPMLicenseFromURL(pkgName, url string) (diligent.Dep, error) {
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

func getNPMLicense(pkgName, version string) (diligent.Dep, error) {
	npmURL := fmt.Sprintf("https://registry.npmjs.org/%s/%s", strings.Replace(url.QueryEscape(pkgName), "%40", "@", 1), url.QueryEscape(version))
	dep, err := getNPMLicenseFromURL(pkgName, npmURL)
	if err == nil {
		return dep, err
	}
	// it seems for scoped packages (e.g. @angular/router) URLs with exact versions
	// like https://registry.npmjs.org/@angular%2Fupgrade/4.4.5 don't work
	// but https://registry.npmjs.org/@angular%2Fupgrade/=4.4.5 do
	// lets try that, if it succeeds great, if not, return the original results
	npmURL = fmt.Sprintf("https://registry.npmjs.org/%s/=%s", strings.Replace(url.QueryEscape(pkgName), "%40", "@", 1), url.QueryEscape(version))
	if dep2, err2 := getNPMLicenseFromURL(pkgName, npmURL); err2 == nil {
		return dep2, err2
	}
	return dep, err
}
