package main

import (
	"fmt"
	"github.com/senseyeio/diligent"
	"log"
)

func isInWhitelist(d diligent.Dep) bool {
	for _, w := range licenseWhitelist {
		if w == d.License.Identifier {
			return true
		}
	}
	return false
}

func whitelistProvided() bool {
	return licenseWhitelist != nil
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
		if whitelistProvided() && isInWhitelist(d) == false {
			return fmt.Errorf("dependency %s has license %s which is not in your license whitelist", d.Name, d.License.Identifier)
		}
	}
	return nil
}
