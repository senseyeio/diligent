package pretty

import (
	"github.com/senseyeio/diligent"
	"io"
	"text/tabwriter"
)

const (
	minColWidth = 5
	tabWidth    = 0 // setting this to 0, as we're not using tabs as padding chars
	padding     = 2
	padChar     = ' '
	flags       = 0
	tab         = "\t"
	newline     = "\n"
)

type pretty struct{}

// NewReporter returns a Reporter which outputs the discovered licenses to a CSV file
func NewReporter() diligent.Reporter {
	return &pretty{}
}

func writeStrings(w io.Writer, strings ...string) error {
	for _, s := range strings {
		_, err := w.Write([]byte(s))

		if err != nil {
			return err
		}
	}

	return nil
}

// Report outputs the dependencies and their licenses in tabulated form to stdout
func (c *pretty) Report(w io.Writer, deps []diligent.Dep) error {
	writer := tabwriter.NewWriter(w, minColWidth, tabWidth, padding, padChar, flags)

	for _, d := range deps {
		err := writeStrings(writer, d.Name, tab, d.License.Name, newline)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}
