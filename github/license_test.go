package github_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/github"
)

func TestIsCompatibleURL(t *testing.T) {
	cases := []struct {
		url        string
		compatible bool
	}{{
		"https://github.com/senseyeio/spaniel",
		true,
	}, {
		"https://senseye.io/senseyeio/spaniel",
		false,
	}, {
		"https://github.com/senseyeio",
		false,
	}, {
		"not-a-url",
		false,
	}}
	target := github.New("https://api.github.com")
	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			compatible := target.IsCompatibleURL(c.url)
			if compatible != c.compatible {
				t.Errorf("expected %t got %t", c.compatible, compatible)
			}
		})
	}
}

func TestGetLicenseFromURL(t *testing.T) {
	cases := []struct {
		d          string
		in         string
		handler    http.HandlerFunc
		expLID     string
		expFailure bool
	}{{
		"should lookup license from github API",
		"https://github.com/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET got %s", r.Method)
			}
			if r.URL.Path != "/repos/senseyeio/spaniel/license" {
				t.Errorf("unexpected path %s", r.URL.Path)
			}
			w.Write([]byte("{\"license\":{\"spdx_id\":\"MIT\"}}"))
			w.WriteHeader(http.StatusOK)
		}),
		"MIT",
		false,
	}, {
		"should fail if not github URL",
		"https://senseye.io/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}),
		"",
		true,
	}, {
		"should fail if github fails",
		"https://github.com/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
		"",
		true,
	}, {
		"should fail if github returns unexpected body",
		"https://github.com/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"license\":{\"noID\":\"it's missing\"}}"))
			w.WriteHeader(http.StatusOK)
		}),
		"",
		true,
	}, {
		"should fail if github returns non json body",
		"https://github.com/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{{"))
			w.WriteHeader(http.StatusOK)
		}),
		"",
		true,
	}, {
		"should fail if github returns an unknown license ID",
		"https://github.com/senseyeio/spaniel",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"license\":{\"spdx_id\":\"woowoo\"}}"))
			w.WriteHeader(http.StatusOK)
		}),
		"",
		true,
	}}

	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			ts := httptest.NewServer(c.handler)
			defer ts.Close()
			target := github.New(ts.URL)
			l, err := target.GetLicenseFromURL(c.in)
			if (err != nil) != c.expFailure {
				t.Errorf("expected failure: %t, got %v", c.expFailure, err)
			}
			if c.expFailure == false {
				expL, _ := diligent.GetLicenseFromIdentifier(c.expLID)
				if expL != l {
					t.Errorf("expected license %+v, got %+v", expL, l)
				}
			}
		})
	}
}
