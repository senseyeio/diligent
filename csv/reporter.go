package csv

import (
	encCSV "encoding/csv"
	"github.com/senseyeio/diligent"
	"io"
)

type csv struct {
}

// NewReporter returns a Reporter which outputs the discovered licenses to a CSV file
func NewReporter() diligent.Reporter {
	return &csv{}
}

// Report outputs the dependencies and their licenses to a CSV file
func (c *csv) Report(w io.Writer, deps []diligent.Dep) error {
	writer := encCSV.NewWriter(w)

	if err := writer.Write([]string{"Name", "License ID", "License Name", "License URL"}); err != nil {
		return err
	}
	for _, d := range deps {
		if err := writer.Write([]string{d.Name, d.License.Identifier, d.License.Name, d.License.URL}); err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}
