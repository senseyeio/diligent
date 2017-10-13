package csv

import (
	encCSV "encoding/csv"
	"github.com/senseyeio/diligent"
	"os"
)

type csv struct {
	filePath string
}

func NewReporter(filePath string) diligent.Reporter {
	return &csv{filePath}
}

func (c *csv) Report(deps []diligent.Dep) error {
	f, err := os.Create(c.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := encCSV.NewWriter(f)

	if err = writer.Write([]string{"Name", "License ID", "License Name", "License URL"}); err != nil {
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
