package stdout

import (
	"fmt"
	"github.com/senseyeio/diligent"
)

type stdout struct{}

func NewReporter() diligent.Reporter {
	return &stdout{}
}

func (s *stdout) Report(deps []diligent.Dep) error {
	for _, dep := range deps {
		fmt.Println(fmt.Sprintf("%s -> %s", dep.Name, dep.License.Name))
	}
	return nil
}
