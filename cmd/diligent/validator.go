package main

import (
	"fmt"

	"github.com/senseyeio/diligent"
	warnpkg "github.com/senseyeio/diligent/warning"
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

func isIgnored(pkgName string) bool {
	for _, i := range ignoreRegex {
		if i.MatchString(pkgName) {
			return true
		}
	}
	return false
}

func ignorePackages(dd []diligent.Dep, ww []diligent.Warning) ([]diligent.Dep, []diligent.Warning) {
	ddOut := make([]diligent.Dep, 0, len(dd))
	for _, d := range dd {
		if !isIgnored(d.Name) {
			ddOut = append(ddOut, d)
		}
	}
	wwOut := make([]diligent.Warning, 0, len(ww))
	for _, w := range ww {
		warn, ok := w.(*warnpkg.Warn)
		if !ok {
			wwOut = append(wwOut, w)
			continue
		}
		if !isIgnored(warn.Dep) {
			wwOut = append(wwOut, w)
		}
	}
	return ddOut, wwOut
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
