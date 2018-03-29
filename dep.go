package diligent

type Dep struct {
	Name    string
	License License
}

// Warning represents an error whilst processing a dependency
// Warnings are not fatal, like an error, but does mean the license associated with a dependency was not found
type Warning interface {
	Warning() string
}

type Deper interface {
	Name() string
	Dependencies(file []byte) ([]Dep, []Warning, error)
	IsCompatible(filename string, fileContents []byte) bool
}
