package dep_test

import (
	"testing"
	"github.com/senseyeio/diligent/dep"
	"github.com/senseyeio/diligent"
	"reflect"
	"errors"
	"github.com/senseyeio/diligent/warning"
)

type licenseGetterResponse struct {
	license diligent.License
	err error
}

type mockLicenseGetter struct {
	responses map[string]licenseGetterResponse
	t *testing.T
}

func newMockLicenseGetter(t *testing.T, responses map[string]licenseGetterResponse) *mockLicenseGetter {
	return &mockLicenseGetter{
		responses: responses,
		t: t,
	}
}

func (mlg *mockLicenseGetter) GetLicense(packagePath string) (diligent.License, error) {
	resp, ok := mlg.responses[packagePath]
	if !ok {
		mlg.t.Errorf("mock not expecting %s", packagePath)
	}
	return resp.license, resp.err
}

func TestName(t *testing.T) {
	mockLG := newMockLicenseGetter(t, map[string]licenseGetterResponse{})
	target := dep.New(mockLG)
	if target.Name() != "dep" {
		t.Error("expected 'dep'")
	}
}

var compatibleTests = []struct {
	in string
	fileContents []byte
	out bool
}{
	{"Gopkg.lock", []byte{}, true},
	{"Gopkg.lock.old", []byte{}, false},
	{"gopkg.lock", []byte{}, false},
	{"Gopkg.toml", []byte{}, false},
	{"package.json", []byte{}, false},
	{"random-Gopkg.lock", []byte{}, false},
}

func TestIsCompatible(t *testing.T) {
	for _, tt := range compatibleTests {
		t.Run(tt.in, func(t *testing.T) {
			mockLG := newMockLicenseGetter(t, map[string]licenseGetterResponse{})
			target := dep.New(mockLG)
			compatible := target.IsCompatible(tt.in, tt.fileContents)
			if compatible != tt.out {
				t.Errorf("got %v, want %v", compatible, tt.out)
			}
		})
	}
}

var depTests = []struct {
	description string
	in []byte
	getLicenseLUT map[string]licenseGetterResponse
	depsOut []diligent.Dep
	warnsOut []diligent.Warning
	errOut bool
}{{
	"single dependency",
	[]byte(`
[[projects]]
  name = "github.com/inconshreveable/mousetrap"
  packages = ["."]
  revision = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75"
  version = "v1.0"
`),
  map[string]licenseGetterResponse{
  	"github.com/inconshreveable/mousetrap": {
  		err: nil,
  		license: diligent.License{Identifier:"MIT"},
	},
  },
  []diligent.Dep{{
  	Name: "github.com/inconshreveable/mousetrap",
  	License: diligent.License{Identifier:"MIT"},
  }},
  []diligent.Warning{},
  false,
}, {
	"multiple dependencues",
	[]byte(`
[[projects]]
  name = "github.com/inconshreveable/mousetrap"
  packages = ["."]
  revision = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75"
  version = "v1.0"

[[projects]]
  name = "github.com/pelletier/go-toml"
  packages = ["."]
  revision = "acdc4509485b587f5e675510c4f2c63e90ff68a8"
  version = "v1.1.0"
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err: nil,
			license: diligent.License{Identifier:"MIT"},
		},
		"github.com/pelletier/go-toml": {
			err: nil,
			license: diligent.License{Identifier:"DOC"},
		},
	},
	[]diligent.Dep{{
		Name: "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier:"MIT"},
	}, {
		Name: "github.com/pelletier/go-toml",
		License: diligent.License{Identifier:"DOC"},
	}},
	[]diligent.Warning{},
	false,
}, {
	"part failure dependencies",
	[]byte(`
[[projects]]
  name = "github.com/inconshreveable/mousetrap"
  packages = ["."]
  revision = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75"
  version = "v1.0"

[[projects]]
  name = "github.com/pelletier/go-toml"
  packages = ["."]
  revision = "acdc4509485b587f5e675510c4f2c63e90ff68a8"
  version = "v1.1.0"
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err: nil,
			license: diligent.License{Identifier:"MIT"},
		},
		"github.com/pelletier/go-toml": {
			err: errors.New("error"),
			license: diligent.License{},
		},
	},
	[]diligent.Dep{{
		Name: "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier:"MIT"},
	}},
	[]diligent.Warning{
		warning.New("github.com/pelletier/go-toml", "error"),
	},
	false,
}, {
	"failure getting any dependencies",
	[]byte(`
[[projects]]
  name = "github.com/inconshreveable/mousetrap"
  packages = ["."]
  revision = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75"
  version = "v1.0"

[[projects]]
  name = "github.com/pelletier/go-toml"
  packages = ["."]
  revision = "acdc4509485b587f5e675510c4f2c63e90ff68a8"
  version = "v1.1.0"
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err: errors.New("eeek"),
			license: diligent.License{},
		},
		"github.com/pelletier/go-toml": {
			err: errors.New("error"),
			license: diligent.License{},
		},
	},
	[]diligent.Dep{},
	[]diligent.Warning{
		warning.New("github.com/inconshreveable/mousetrap", "eeek"),
		warning.New("github.com/pelletier/go-toml", "error"),
	},
	false,
}, {
	"toml parsing failure",
	[]byte(`
{
	"woops": "it is json"
}
`),
	map[string]licenseGetterResponse{},
	[]diligent.Dep{},
	[]diligent.Warning{},
	true,
}}

func TestDependencies(t *testing.T) {
	for _, tt := range depTests {
		t.Run(tt.description, func(t *testing.T) {
			mockLG := newMockLicenseGetter(t, tt.getLicenseLUT)
			target := dep.New(mockLG)
			d, w, e := target.Dependencies(tt.in)
			if (len(d) > 0 || len(tt.depsOut) > 0) && reflect.DeepEqual(d, tt.depsOut) == false {
				t.Errorf("deps: got %v, want %v", d, tt.depsOut)
			}
			if (len(d) > 0 || len(tt.depsOut) > 0) && reflect.DeepEqual(w, tt.warnsOut) == false {
				t.Errorf("warnings: got %v, want %v", w, tt.warnsOut)
			}
			isErr := e != nil
			if tt.errOut != isErr {
				t.Errorf("error: got %v, want %v", isErr, tt.errOut)
			}
		})
	}
}