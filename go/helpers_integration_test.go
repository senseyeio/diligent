// +build integration

package _go_test

import (
	"testing"

	"errors"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/go"
)

type webLicenseGetterResponse struct {
	isCompatible bool
	license      diligent.License
	err          error
}

type mockWebLicenseGetter struct {
	responses map[string]webLicenseGetterResponse
}

func (m mockWebLicenseGetter) IsCompatibleURL(s string) bool {
	r, ok := m.responses[s]
	return ok && r.isCompatible
}
func (m mockWebLicenseGetter) GetLicenseFromURL(s string) (diligent.License, error) {
	r, ok := m.responses[s]
	if !ok {
		return diligent.License{}, errors.New("not mocked")
	}
	return r.license, r.err
}

func TestWebLicenseGetter(t *testing.T) {
	cases := []struct {
		d          string
		pkgInput   string
		expFailure bool
	}{{
		"should use license getter",
		"test/one",
		false,
	}, {
		"should support deep packages",
		"test/one/two/three/four/five",
		false,
	}, {
		"should handle invalid package identifiers",
		"test",
		true,
	}}
	mock := mockWebLicenseGetter{map[string]webLicenseGetterResponse{
		"https://test/one": {
			isCompatible: true,
			license:      diligent.License{Name: "test-license"},
			err:          nil,
		},
		"https://test/one/two": {
			isCompatible: false,
		},
	}}
	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			target := _go.NewLicenseGetter(mock)
			l, err := target.GetLicense(c.pkgInput)
			if (err != nil) != c.expFailure {
				t.Errorf("expected failure: %t, got %v", c.expFailure, err)
			}
			if c.expFailure == false && l.Name != "test-license" {
				t.Errorf("expected license test-license, got %+v", l)
			}
		})
	}
}

func TestGoGetLicenseLookup(t *testing.T) {
	mock := mockWebLicenseGetter{map[string]webLicenseGetterResponse{}}
	target := _go.NewLicenseGetter(mock)

	t.Run("should fail with an invalid package", func(t *testing.T) {
		_, err := target.GetLicense("not/a/real/package")
		if err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("should succeed with an valid package", func(t *testing.T) {
		l, err := target.GetLicense("github.com/senseyeio/spaniel")
		if err != nil {
			t.Errorf("did not expect an error, got %v", err)
		}
		if l.Identifier != "MIT" {
			t.Errorf("expected MIT license, got %+v", l)
		}
	})
}
