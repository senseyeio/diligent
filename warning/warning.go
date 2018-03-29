package warning

import "github.com/senseyeio/diligent"

// New returns a Warning. It includes the name of the dependency and a message describing the problem.
func New(dependency string, message string) diligent.Warning {
	return &warn{
		msg: message,
		dep: dependency,
	}
}

type warn struct {
	msg string
	dep string
}

// Warning implements diligent.Warning
func (w *warn) Warning() string {
	return "Failed to determine license for " + w.dep + ": " + w.msg
}
