package main

import (
	"fmt"
	"github.com/senseyeio/diligent"
	"log"
)

func isInWhitelist(l diligent.License) bool {
	if len(licenseWhitelist) == 0 {
		return true
	}
	for _, w := range licenseWhitelist {
		if w == l.Identifier {
			return true
		}
	}
	return false
}

func checkWhitelist() error {
	for _, w := range licenseWhitelist {
		l := diligent.GetLicenseFromIdentifier(w)
		if diligent.IsUnknownLicense(l) {
			log.Printf("Whitelisted license %s is not a known license identifier", w)
		}
	}
	return nil
}

func validateDependencies(deps []diligent.Dep) error {
	for _, d := range deps {
		if isInWhitelist(d.License) == false {
			return fmt.Errorf("dependency %s has license %s which is not in your license whitelist", d.Name, d.License.Identifier)
		}
	}
	return nil
}
