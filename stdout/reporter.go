package stdout

import (
	"fmt"

	"github.com/senseyeio/diligent"
)

type stdout struct{}

// NewReporter creates a stdout Reporter
// Each dependency's license will be written to stdout
func NewReporter() diligent.Reporter {
	return &stdout{}
}

// Report outputs the dependencies and their licenses to stdout
func (s *stdout) Report(deps []diligent.Dep) error {
	for _, dep := range deps {
		fmt.Println(fmt.Sprintf("%s -> %s", dep.Name, dep.License.Name))
	}
	return nil
}
