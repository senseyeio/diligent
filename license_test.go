package diligent_test

import (
	"testing"

	"github.com/senseyeio/diligent"
)

func testForPropFreeLicenses(t *testing.T, identifierMap map[string]bool) {
	if _, ok := identifierMap["SISSL"]; !ok {
		t.Errorf("expected SISSL, got %v", identifierMap)
	}
	if _, ok := identifierMap["SISSL-1.2"]; !ok {
		t.Errorf("expected SISSL-1.2, got %v", identifierMap)
	}
	if _, ok := identifierMap["Watcom-1.0"]; !ok {
		t.Errorf("expected Watcom-1.0, got %v", identifierMap)
	}
	if _, ok := identifierMap["Unicode-TOU"]; !ok {
		t.Errorf("expected Unicode-TOU, got %v", identifierMap)
	}
}

func TestReplaceCategoriesWithIdentifiers(t *testing.T) {
	out := diligent.ReplaceCategoriesWithIdentifiers([]string{"MIT", "proprietary-free"})
	if len(out) != 5 {
		t.Errorf("expected five identifiers, got %v", len(out))
	}
	m := map[string]bool{}
	for _, o := range out {
		m[o] = true
	}
	testForPropFreeLicenses(t, m)
	if _, ok := m["MIT"]; !ok {
		t.Errorf("expected MIT, got %v", out)
	}
}

func TestGetCategoryLicenses(t *testing.T) {
	out := diligent.GetCategoryLicenses(diligent.ProprietaryFree)
	if len(out) != 4 {
		t.Errorf("expected four identifiers, got %v", len(out))
	}
	m := map[string]bool{}
	for _, o := range out {
		m[o.Identifier] = true
	}
	testForPropFreeLicenses(t, m)
}

func TestGetLicenses(t *testing.T) {
	out := diligent.GetLicenses()
	m := map[string]bool{}
	for _, o := range out {
		m[o.Identifier] = true
	}
	testForPropFreeLicenses(t, m)
	if _, ok := m["MIT"]; !ok {
		t.Errorf("expected MIT, got %v", out)
	}
}

func TestGetLicenseFromIdentifier(t *testing.T) {
	cases := []struct {
		d             string
		in            string
		outIdentifier string
		expFailure    bool
	}{
		{"standard identifier", "MIT", "MIT", false},
		{"unknown identifier", "woowoo", "", true},
		{"empty identifier", "", "", true},
		{"handle non standard 'NewBSD'", "NewBSD", "BSD-3-Clause", false},
		{"handle non standard 'FreeBSD'", "FreeBSD", "BSD-2-Clause", false},
	}

	for _, c := range cases {
		t.Run(c.d, func(t *testing.T) {
			out, err := diligent.GetLicenseFromIdentifier(c.in)
			if (err != nil) != c.expFailure {
				t.Errorf("expecting error: %t, got %s", c.expFailure, err)
			}
			if !c.expFailure && out.Identifier != c.outIdentifier {
				t.Errorf("expecting %s, got %+v", c.outIdentifier, out)
			}
		})
	}
}
