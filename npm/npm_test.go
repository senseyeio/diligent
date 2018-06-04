package npm_test

import (
	"testing"

	"net/http"
	"reflect"

	"net/http/httptest"

	"sort"

	"net/url"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/npm"
	"github.com/senseyeio/diligent/warning"
)

func TestName(t *testing.T) {
	target := npm.New("")
	if target.Name() != "npm" {
		t.Error("expected 'npm'")
	}
}

func TestIsCompatible(t *testing.T) {
	var cases = []struct {
		in           string
		fileContents []byte
		out          bool
	}{
		{"package.json", []byte{}, true},
		{"package.json.new", []byte{}, false},
		{"Package.json", []byte{}, false},
		{"package.lock", []byte{}, false},
		{"vendor.json", []byte{}, false},
		{"random-package.json", []byte{}, false},
	}

	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			target := npm.New("")
			compatible := target.IsCompatible(tt.in, tt.fileContents)
			if compatible != tt.out {
				t.Errorf("got %v, want %v", compatible, tt.out)
			}
		})
	}
}

func TestDependencies(t *testing.T) {
	pathAndQuery := func(url *url.URL) string {
		return url.Path + "?" + url.RawQuery
	}
	cases := []struct {
		description string
		config      npm.Config
		in          []byte
		handler     http.HandlerFunc
		depsOut     map[string]string
		warnsOut    []diligent.Warning
		errOut      bool
	}{{
		"should handle only dependencies by default",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0"
				},
				"devDependencies": {
					"cypress": "2.1.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET got %s", r.Method)
			}
			if pathAndQuery(r.URL) != "/d3?version=5.0.0" {
				t.Errorf("unexpected path %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"license\":\"MIT\"}"))
		}),
		map[string]string{
			"d3": "MIT",
		},
		[]diligent.Warning{},
		false,
	}, {
		"should be able to handle multiple dependencies",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "^5.0.0",
					"cypress": "2.1.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch pathAndQuery(r.URL) {
			case "/d3?version=%5E5.0.0":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"license\":\"GPL-3.0\"}"))
			case "/cypress?version=2.1.0":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"license\":\"MIT\"}"))
			default:
				t.Errorf("unexpected path %s", pathAndQuery(r.URL))
			}
		}),
		map[string]string{
			"d3":      "GPL-3.0",
			"cypress": "MIT",
		},
		[]diligent.Warning{},
		false,
	}, {
		"should support part failures",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "~5.0.0",
					"cypress": "2.1.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch pathAndQuery(r.URL) {
			case "/d3?version=~5.0.0":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"license\":\"GPL-3.0\"}"))
			case "/cypress?version=2.1.0":
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"failed\"}"))
			default:
				t.Errorf("unexpected path %s", pathAndQuery(r.URL))
			}
		}),
		map[string]string{
			"d3": "GPL-3.0",
		},
		[]diligent.Warning{
			warning.New("cypress", "requested failed with status 500"),
		},
		false,
	}, {
		"should support all deps failing",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0",
					"cypress": "2.1.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch pathAndQuery(r.URL) {
			case "/d3?version=5.0.0":
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"failed\"}"))
			case "/cypress?version=2.1.0":
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{\"error\":\"failed\"}"))
			default:
				t.Errorf("unexpected path %s", pathAndQuery(r.URL))
			}
		}),
		map[string]string{},
		[]diligent.Warning{
			warning.New("d3", "requested failed with status 500"),
			warning.New("cypress", "requested failed with status 500")},
		false,
	}, {
		"should be capable of including devDependencies",
		npm.Config{DevDependencies: true},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0"
				},
				"devDependencies": {
					"cypress": "2.1.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET got %s", r.Method)
			}
			if pathAndQuery(r.URL) != "/d3?version=5.0.0" && pathAndQuery(r.URL) != "/cypress?version=2.1.0" {
				t.Errorf("unexpected path %s", pathAndQuery(r.URL))
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"license\":\"MIT\"}"))
		}),
		map[string]string{
			"d3":      "MIT",
			"cypress": "MIT",
		},
		[]diligent.Warning{},
		false,
	}, {
		"package.json parse failure",
		npm.Config{},
		[]byte(`{{`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		map[string]string{},
		[]diligent.Warning{},
		true,
	}, {
		"should fail with unknown license",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"license\":\"woowoo\"}"))
		}),
		map[string]string{},
		[]diligent.Warning{
			warning.New("d3", "license identifier woowoo is not known to diligent"),
		},
		false,
	}, {
		"should fail if response is not valid JSON",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{{"))
		}),
		map[string]string{},
		[]diligent.Warning{
			warning.New("d3", "parsing NPM response failed - invalid JSON"),
		},
		false,
	}, {
		"should warn if no license info in response",
		npm.Config{},
		[]byte(`
			{
				"dependencies": {
					"d3": "5.0.0"
				}
			}
		`),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}),
		map[string]string{},
		[]diligent.Warning{
			warning.New("d3", "no license information in NPM"),
		},
		false,
	}}

	for _, tt := range cases {
		t.Run(tt.description, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()
			target := npm.NewWithOptions(ts.URL, tt.config)
			d, w, e := target.Dependencies(tt.in)
			expectedDeps := make([]diligent.Dep, 0, len(tt.depsOut))
			for depID, lID := range tt.depsOut {
				l, _ := diligent.GetLicenseFromIdentifier(lID)
				expectedDeps = append(expectedDeps, diligent.Dep{depID, l})
			}
			if len(d) > 0 || len(expectedDeps) > 0 {
				sort.Sort(diligent.DepsByName(d))
				sort.Sort(diligent.DepsByName(expectedDeps))
				if reflect.DeepEqual(d, expectedDeps) == false {
					t.Errorf("deps: got %+v, want %+v", d, expectedDeps)
				}
			}
			if (len(w) > 0 || len(tt.warnsOut) > 0) && reflect.DeepEqual(w, tt.warnsOut) == false {
				t.Errorf("warnings: got %+v, want %+v", w, tt.warnsOut)
			}
			isErr := e != nil
			if tt.errOut != isErr {
				t.Errorf("error: got %v, want %v", isErr, tt.errOut)
			}
		})
	}
}
