package warning

import "github.com/senseyeio/diligent"

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

func (w *warn) Warning () string {
	return "Failed to determine license for " + w.dep + ": " + w.msg
}