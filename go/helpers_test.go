package _go_test

import (
	"testing"

	"github.com/senseyeio/diligent/go"
)

func TestInvalidPath(t *testing.T) {
	target := _go.NewLicenseGetter(nil)
	_, err := target.GetLicense("no-components")
	if err == nil {
		t.Error("expected an error")
	}
}
