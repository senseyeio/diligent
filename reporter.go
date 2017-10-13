package diligent

// Reporter takes an array of dependencies and outputs them to a certain medium
type Reporter interface {
	Report(deps []Dep) error
}
