package diligent

import "io"

// Reporter takes an array of dependencies and outputs them to a certain medium
type Reporter interface {
	Report(w io.Writer, deps []Dep) error
}
