package main

import (
	"fmt"

	"github.com/senseyeio/diligent"
)

func isInWhitelist(l diligent.License) bool {
	for _, w := range licenseWhitelist {
		if w == l.Identifier {
			return true
		}
	}
	return false
}

func checkWhitelist() error {
	for _, w := range licenseWhitelist {
		_, err := diligent.GetLicenseFromIdentifier(w)
		if err != nil {
			return fmt.Errorf("whitelisted license '%s' is not a known license identifier", w)
		}
	}
	return nil
}

func validateDependencies(deps []diligent.Dep) []error {
	ee := make([]error, 0, len(deps))
	for _, d := range deps {
		if isInWhitelist(d.License) == false {
			ee = append(ee, fmt.Errorf("dependency '%s' has license '%s' which is not in your license whitelist", d.Name, d.License.Identifier))
		}
	}
	return ee
}
