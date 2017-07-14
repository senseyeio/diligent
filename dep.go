package diligent

type Dep struct {
	Name    string
	License License
}

type Deper interface {
	Name() string
	Dependencies(file []byte, options map[string]interface{}) ([]Dep, error)
	IsCompatible(filename string, fileContents []byte) bool
}
