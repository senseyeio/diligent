package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"errors"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/warning"
)

type packageJSON struct {
	Deps    map[string]string `json:"dependencies"`
	DevDeps map[string]string `json:"devDependencies"`
}

type license string

type repo string
type repoObject struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// WebLicenseGetter retrieves license information from an online source
type WebLicenseGetter interface {
	IsCompatibleURL(s string) bool
	GetLicenseFromURL(s string) (diligent.License, error)
}

func (r *repo) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err == nil {
		if strings.HasPrefix(str, "github:") {
			str = strings.Replace(str, "github:", "git+https://github.com/", 1)
		}
		*r = repo(str)
		return nil
	}
	var ro repoObject
	err = json.Unmarshal(data, &ro)
	if err == nil {
		*r = repo(ro.URL)
		return nil
	}
	return errors.New("unknown format for 'repository'")

}

type npmPackage struct {
	License    *license `json:"license"`
	Repository *repo    `json:"repository"`
}

type npmDeper struct {
	config Config
	url    string
	webLG  WebLicenseGetter
}

// Config allows default options to be altered
type Config struct {
	// DevDependencies can be set to true if you want to gather the licenses of your devDependencies as well as your dependencies
	DevDependencies bool
}

// New returns a Deper capable of dealing with package.json manifest files
func New(url string, webLG WebLicenseGetter) diligent.Deper {
	return NewWithOptions(url, webLG, Config{})
}

// NewWithOptions is identical to New but allows the default options to be overridden
func NewWithOptions(url string, webLG WebLicenseGetter, c Config) diligent.Deper {
	return &npmDeper{c, url, webLG}
}

// Name returns "npm"
func (n *npmDeper) Name() string {
	return "npm"
}

func mergeMaps(to map[string]string, from map[string]string) {
	for pkg, version := range from {
		to[pkg] = version
	}
}

// Dependencies returns the licenses associated with the NPM dependencies
func (n *npmDeper) Dependencies(file []byte) ([]diligent.Dep, []diligent.Warning, error) {
	var pkg packageJSON
	err := json.Unmarshal(file, &pkg)
	if err != nil {
		return nil, nil, err
	}

	licensesToGet := map[string]string{}
	mergeMaps(licensesToGet, pkg.Deps)
	if n.config.DevDependencies {
		mergeMaps(licensesToGet, pkg.DevDeps)
	}

	deps := make([]diligent.Dep, 0, len(licensesToGet))
	warns := make([]diligent.Warning, 0, len(licensesToGet))
	for pkg, version := range licensesToGet {
		l, err := n.getNPMLicense(pkg, version)
		if err != nil {
			warns = append(warns, warning.New(pkg, err.Error()))
		} else {
			deps = append(deps, l)
		}
	}
	return deps, warns, nil
}

// IsCompatible returns true if the filename is package.json
func (n *npmDeper) IsCompatible(filename string) bool {
	return filename == "package.json"
}

func (n *npmDeper) getNPMLicenseFromURL(pkgName, url string) (diligent.Dep, error) {
	resp, err := http.Get(url)
	if err != nil {
		return diligent.Dep{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return diligent.Dep{}, fmt.Errorf("requested failed with status %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return diligent.Dep{}, err
	}

	var packageInfo npmPackage
	err = json.Unmarshal(body, &packageInfo)
	if err != nil {
		return diligent.Dep{}, errors.New("parsing NPM response failed - invalid JSON")
	}

	if packageInfo.License != nil {
		l, err := diligent.GetLicenseFromIdentifier(string(*packageInfo.License))
		if err == nil {
			return diligent.Dep{
				Name:    pkgName,
				License: l,
			}, nil
		}
	}

	if packageInfo.Repository != nil {
		if n.webLG != nil && n.webLG.IsCompatibleURL(string(*packageInfo.Repository)) {
			gitUrl := string(*packageInfo.Repository)
			repoURL := strings.Replace(gitUrl, ".git", "", 1)
			l, err := n.webLG.GetLicenseFromURL(repoURL)
			if err == nil {
				return diligent.Dep{
					Name:    pkgName,
					License: l,
				}, nil
			}
		}
		if strings.HasPrefix(string(*packageInfo.Repository), "git") {
			l, err := diligent.GetLicenseForGit(string(*packageInfo.Repository))
			if err == nil {
				return diligent.Dep{
					Name:    pkgName,
					License: l,
				}, nil
			}
		}
	}

	return diligent.Dep{}, errors.New("no license information in NPM")
}

func (n *npmDeper) getNPMLicense(pkgName, version string) (diligent.Dep, error) {
	npmURL := fmt.Sprintf("%s/%s?version=%s", n.url, strings.Replace(url.QueryEscape(pkgName), "%40", "@", 1), url.QueryEscape(version))
	return n.getNPMLicenseFromURL(pkgName, npmURL)
}

func (l *license) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch raw := raw.(type) {
	case string:
		*l = license(raw)
	case map[string]interface{}:
		if t, ok := raw["type"].(string); ok {
			*l = license(t)
		}
	}
	return nil
}
