package diligent

import "fmt"

// Dep contains a dependency identified by name along with its License information
type Dep struct {
	Name    string
	License License
}

// Warning represents an error whilst processing a dependency
// Warnings are not fatal, like an error, but does mean the license associated with a dependency was not found
type Warning interface {
	Warning() string
}

// Deper is the interface for extracting licenses from manifest files.
// Implementations should interrogate a package manager's manifest file, determine the dependencies and identify their licenses
type Deper interface {
	// Name returns the name of the Deper
	Name() string
	// Dependencies interrogates the manifest file and returns the licenses associated with each dependency
	// If a single dependency cannot be processed, a warning should be returned
	// If no dependencies can be processed, an error should be returned
	Dependencies(file []byte) ([]Dep, []Warning, error)
	// IsCompatible should return true if the Deper can handle the provided manifest file
	IsCompatible(filename string) bool
}

type DepsByName []Dep

func (d DepsByName) Len() int           { return len(d) }
func (d DepsByName) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DepsByName) Less(i, j int) bool { return d[i].Name < d[j].Name }

type Warnings []Warning

func (d Warnings) Len() int           { return len(d) }
func (d Warnings) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d Warnings) Less(i, j int) bool { return d[i].Warning() < d[j].Warning() }

type DepsByLicense []Dep

func (d DepsByLicense) Len() int      { return len(d) }
func (d DepsByLicense) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d DepsByLicense) Less(i, j int) bool {
	if d[i].License.Name == d[j].License.Name {
		return d[i].Name < d[j].Name
	}

	return d[i].License.Name < d[j].License.Name
}

type Deps []Dep

// Dedupe removes duplicate dependencies in place
func (dd Deps) Dedupe() Deps {
	out := make([]Dep, 0, len(dd))
	found := map[string]bool{}
	for _, d := range dd {
		key := fmt.Sprintf("%s-%s", d.Name, d.License.Identifier)
		if _, ok := found[key]; !ok {
			out = append(out, d)
			found[key] = true
		}
	}
	return out
}
