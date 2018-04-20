package govendor_test

import (
"testing"
"github.com/senseyeio/diligent"
"reflect"
"errors"
"github.com/senseyeio/diligent/warning"
	"github.com/senseyeio/diligent/govendor"
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
	target := govendor.New(mockLG)
	if target.Name() != "govendor" {
		t.Error("expected 'govendor'")
	}
}

var compatibleTests = []struct {
	in string
	fileContents []byte
	out bool
}{
	{"vendor.json", []byte{}, true},
	{"vendor.json.new", []byte{}, false},
	{"Vendor.json", []byte{}, false},
	{"vendor.lock", []byte{}, false},
	{"package.json", []byte{}, false},
	{"random-vendor.json", []byte{}, false},
}

func TestIsCompatible(t *testing.T) {
	for _, tt := range compatibleTests {
		t.Run(tt.in, func(t *testing.T) {
			mockLG := newMockLicenseGetter(t, map[string]licenseGetterResponse{})
			target := govendor.New(mockLG)
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
{
	"package": [
		{
			"checksumSHA1": "KxX/Drph+byPXBFIXaCZaCOAnrU=",
			"path": "github.com/go-logfmt/logfmt",
			"revision": "390ab7935ee28ec6b286364bba9b4dd6410cb3d5",
			"revisionTime": "2016-11-15T14:25:13Z"
		}
	]
}
`),
	map[string]licenseGetterResponse{
		"github.com/go-logfmt/logfmt": {
			err: nil,
			license: diligent.License{Identifier:"MIT"},
		},
	},
	[]diligent.Dep{{
		Name: "github.com/go-logfmt/logfmt",
		License: diligent.License{Identifier:"MIT"},
	}},
	[]diligent.Warning{},
	false,
}, {
	"multiple dependencies",
	[]byte(`
{
	"package": [
		{
			"checksumSHA1": "KxX/Drph+byPXBFIXaCZaCOAnrU=",
			"path": "github.com/go-logfmt/logfmt",
			"revision": "390ab7935ee28ec6b286364bba9b4dd6410cb3d5",
			"revisionTime": "2016-11-15T14:25:13Z"
		},
		{
			"checksumSHA1": "j6vhe49MX+dyHR9rU91P6vMx55o=",
			"path": "github.com/go-stack/stack",
			"revision": "817915b46b97fd7bb80e8ab6b69f01a53ac3eebf",
			"revisionTime": "2017-07-24T01:23:01Z"
		}
	]
}
`),
	map[string]licenseGetterResponse{
		"github.com/go-logfmt/logfmt": {
			err: nil,
			license: diligent.License{Identifier:"MIT"},
		},
		"github.com/go-stack/stack": {
			err: nil,
			license: diligent.License{Identifier:"DOC"},
		},
	},
	[]diligent.Dep{{
		Name: "github.com/go-logfmt/logfmt",
		License: diligent.License{Identifier:"MIT"},
	}, {
		Name: "github.com/go-stack/stack",
		License: diligent.License{Identifier:"DOC"},
	}},
	[]diligent.Warning{},
	false,
}, {
	"part failure dependencies",
	[]byte(`
{
	"package": [
		{
			"checksumSHA1": "KxX/Drph+byPXBFIXaCZaCOAnrU=",
			"path": "github.com/go-logfmt/logfmt",
			"revision": "390ab7935ee28ec6b286364bba9b4dd6410cb3d5",
			"revisionTime": "2016-11-15T14:25:13Z"
		},
		{
			"checksumSHA1": "j6vhe49MX+dyHR9rU91P6vMx55o=",
			"path": "github.com/go-stack/stack",
			"revision": "817915b46b97fd7bb80e8ab6b69f01a53ac3eebf",
			"revisionTime": "2017-07-24T01:23:01Z"
		}
	]
}
`),
	map[string]licenseGetterResponse{
		"github.com/go-logfmt/logfmt": {
			err: nil,
			license: diligent.License{Identifier:"MIT"},
		},
		"github.com/go-stack/stack": {
			err: errors.New("error"),
			license: diligent.License{},
		},
	},
	[]diligent.Dep{{
		Name: "github.com/go-logfmt/logfmt",
		License: diligent.License{Identifier:"MIT"},
	}},
	[]diligent.Warning{
		warning.New("github.com/go-stack/stack", "error"),
	},
	false,
}, {
	"failure getting any dependencies",
	[]byte(`
{
	"package": [
		{
			"checksumSHA1": "KxX/Drph+byPXBFIXaCZaCOAnrU=",
			"path": "github.com/go-logfmt/logfmt",
			"revision": "390ab7935ee28ec6b286364bba9b4dd6410cb3d5",
			"revisionTime": "2016-11-15T14:25:13Z"
		},
		{
			"checksumSHA1": "j6vhe49MX+dyHR9rU91P6vMx55o=",
			"path": "github.com/go-stack/stack",
			"revision": "817915b46b97fd7bb80e8ab6b69f01a53ac3eebf",
			"revisionTime": "2017-07-24T01:23:01Z"
		}
	]
}
`),
	map[string]licenseGetterResponse{
		"github.com/go-logfmt/logfmt": {
			err: errors.New("eeek"),
			license: diligent.License{},
		},
		"github.com/go-stack/stack": {
			err: errors.New("error"),
			license: diligent.License{},
		},
	},
	[]diligent.Dep{},
	[]diligent.Warning{
		warning.New("github.com/go-logfmt/logfmt", "eeek"),
		warning.New("github.com/go-stack/stack", "error"),
	},
	false,
}, {
	"json parsing failure",
	[]byte(`
{
	"not valid json""
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
			target := govendor.New(mockLG)
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