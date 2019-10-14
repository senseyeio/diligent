package gomod_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/senseyeio/diligent/gomod"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/warning"
)

type licenseGetterResponse struct {
	license diligent.License
	err     error
}

type mockLicenseGetter struct {
	responses map[string]licenseGetterResponse
	t         *testing.T
}

func newMockLicenseGetter(t *testing.T, responses map[string]licenseGetterResponse) *mockLicenseGetter {
	return &mockLicenseGetter{
		responses: responses,
		t:         t,
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
	target := gomod.New(mockLG)
	if target.Name() != "gomod" {
		t.Error("expected 'gomod'")
	}
}

var compatibleTests = []struct {
	in  string
	out bool
}{
	{"go.mod", true},
	{"go.mod.old", false},
	{"go.sum", false},
	{"Gopkg.toml", false},
	{"package.json", false},
	{"random-Gopkg.lock", false},
}

func TestIsCompatible(t *testing.T) {
	for _, tt := range compatibleTests {
		t.Run(tt.in, func(t *testing.T) {
			mockLG := newMockLicenseGetter(t, map[string]licenseGetterResponse{})
			target := gomod.New(mockLG)
			compatible := target.IsCompatible(tt.in)
			if compatible != tt.out {
				t.Errorf("got %v, want %v", compatible, tt.out)
			}
		})
	}
}

var depTests = []struct {
	description   string
	in            []byte
	getLicenseLUT map[string]licenseGetterResponse
	depsOut       []diligent.Dep
	warnsOut      []diligent.Warning
	errOut        bool
}{{
	"single dependency",
	[]byte(`
module my/thing
require github.com/inconshreveable/mousetrap v1.0.0
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err:     nil,
			license: diligent.License{Identifier: "MIT"},
		},
	},
	[]diligent.Dep{{
		Name:    "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier: "MIT"},
	}},
	[]diligent.Warning{},
	false,
}, {
	"multiple dependencies",
	[]byte(`
module my/thing
require (
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/pelletier/go-toml v1.1.0
)
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err:     nil,
			license: diligent.License{Identifier: "MIT"},
		},
		"github.com/pelletier/go-toml": {
			err:     nil,
			license: diligent.License{Identifier: "DOC"},
		},
	},
	[]diligent.Dep{{
		Name:    "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier: "MIT"},
	}, {
		Name:    "github.com/pelletier/go-toml",
		License: diligent.License{Identifier: "DOC"},
	}},
	[]diligent.Warning{},
	false,
}, {
	"part failure dependencies",
	[]byte(`
module my/thing
require (
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/pelletier/go-toml v1.1.0
)
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err:     nil,
			license: diligent.License{Identifier: "MIT"},
		},
		"github.com/pelletier/go-toml": {
			err:     errors.New("error"),
			license: diligent.License{},
		},
	},
	[]diligent.Dep{{
		Name:    "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier: "MIT"},
	}},
	[]diligent.Warning{
		warning.New("github.com/pelletier/go-toml", "error"),
	},
	false,
}, {
	"failure getting any dependencies",
	[]byte(`
module my/thing
require (
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/pelletier/go-toml v1.1.0
)
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err:     errors.New("eeek"),
			license: diligent.License{},
		},
		"github.com/pelletier/go-toml": {
			err:     errors.New("error"),
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
	"parsing failure",
	[]byte(`
{
	"woops": "wrong format"
}
`),
	map[string]licenseGetterResponse{},
	[]diligent.Dep{},
	[]diligent.Warning{},
	true,
}, {
	"replacements",
	[]byte(`
module my/thing
require (
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/pelletier/go-toml v1.1.0
)
replace github.com/pelletier/go-toml v1.1.0 => github.com/russross/blackfriday/v2 v2.0.1
`),
	map[string]licenseGetterResponse{
		"github.com/inconshreveable/mousetrap": {
			err:     nil,
			license: diligent.License{Identifier: "MIT"},
		},
		"github.com/russross/blackfriday/v2": {
			err:     nil,
			license: diligent.License{Identifier: "REP"},
		},
	},
	[]diligent.Dep{{
		Name:    "github.com/inconshreveable/mousetrap",
		License: diligent.License{Identifier: "MIT"},
	}, {
		Name:    "github.com/russross/blackfriday/v2",
		License: diligent.License{Identifier: "REP"},
	}},
	[]diligent.Warning{},
	false,
}}

func TestDependencies(t *testing.T) {
	for _, tt := range depTests {
		t.Run(tt.description, func(t *testing.T) {
			mockLG := newMockLicenseGetter(t, tt.getLicenseLUT)
			target := gomod.New(mockLG)
			d, w, e := target.Dependencies(tt.in)
			if (len(d) > 0 || len(tt.depsOut) > 0) && reflect.DeepEqual(d, tt.depsOut) == false {
				t.Errorf("deps: got %v, want %v", d, tt.depsOut)
			}
			if (len(w) > 0 || len(tt.warnsOut) > 0) && reflect.DeepEqual(w, tt.warnsOut) == false {
				t.Errorf("warnings: got %v, want %v", w, tt.warnsOut)
			}
			isErr := e != nil
			if tt.errOut != isErr {
				t.Errorf("error: got %v, want %v", isErr, tt.errOut)
			}
		})
	}
}
