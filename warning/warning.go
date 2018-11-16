package warning

import "github.com/senseyeio/diligent"

// New returns a Warning. It includes the name of the dependency and a message describing the problem.
func New(dependency string, message string) diligent.Warning {
	return &Warn{
		Msg: message,
		Dep: dependency,
	}
}

type Warn struct {
	Msg string
	Dep string
}

// Warning implements diligent.Warning
func (w *Warn) Warning() string {
	return "Failed to determine license for " + w.Dep + ": " + w.Msg
}
