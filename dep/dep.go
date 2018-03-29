package dep

import (
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/go"
	"github.com/senseyeio/diligent/warning"
)

type lockedProject struct {
	Name string `toml:"name"`
}

type lock struct {
	Projects []lockedProject `toml:"projects"`
}

type dep struct{}

// New returns a Deper capable of handling dep manifest files
func New() diligent.Deper {
	return &dep{}
}

// Name returns "dep"
func (d *dep) Name() string {
	return "dep"
}

// Dependencies returns the licenses of the go packages defined within the dep manifest
func (d *dep) Dependencies(file []byte) ([]diligent.Dep, []diligent.Warning, error) {
	var l lock
	err := toml.Unmarshal(file, &l)
	if err != nil {
		return nil, nil, err
	}

	deps := make([]diligent.Dep, 0, len(l.Projects))
	warns := make([]diligent.Warning, 0, len(l.Projects))
	for _, pkg := range l.Projects {
		l, err := _go.GetLicense(pkg.Name)
		if err != nil {
			warns = append(warns, warning.New(pkg.Name, err.Error()))
		} else {
			deps = append(deps, diligent.Dep{
				Name:    pkg.Name,
				License: l,
			})
		}
	}
	return deps, warns, nil
}

// IsCompatible returns true if the filename is Gopkg.lock
func (d *dep) IsCompatible(filename string, fileContents []byte) bool {
	return strings.Index(filename, "Gopkg.lock") != -1
}
