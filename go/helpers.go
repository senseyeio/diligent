package _go

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"strings"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
	"github.com/senseyeio/diligent"
)

// LicenseGetter provides methods to retrieve the licenses associated with go packages
type LicenseGetter struct {
	webLG WebLicenseGetter
}

// NewLicenseGetter returns a new instance of LicenseGetter using the provided WebLicenseGetter where possible
func NewLicenseGetter(webLG WebLicenseGetter) *LicenseGetter {
	return &LicenseGetter{webLG}
}

// WebLicenseGetter retrieves license information from an online source
type WebLicenseGetter interface {
	IsCompatibleURL(s string) bool
	GetLicenseFromURL(s string) (diligent.License, error)
}

func goPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

// GetLicense will return the license associated with a given go package
func (lg *LicenseGetter) GetLicense(packagePath string) (diligent.License, error) {
	components := strings.Split(packagePath, "/")
	// in some go vendoring solutions full paths to packages are defined as dependencies
	// need to look for the base package identifier so github.com/aws/aws-sdk-go/aws becomes github.com/aws/aws-sdk-go
	if len(components) < 2 {
		return diligent.License{}, errors.New("invalid go package path")
	}
	// try a three component base package, if possible, as it is most common
	if len(components) >= 3 {
		l, err := lg.getLicenseForBasePackage(strings.Join(components[:3], "/"))
		if err == nil {
			return l, nil
		}
	}
	// can have libraries with just two components, for example gopkg.in/mgo.v2
	return lg.getLicenseForBasePackage(strings.Join(components[:2], "/"))
}

func (lg *LicenseGetter) getLicenseForBasePackage(pkg string) (diligent.License, error) {
	if lg.webLG.IsCompatibleURL(fmt.Sprintf("https://%s", pkg)) {
		l, err := lg.webLG.GetLicenseFromURL(fmt.Sprintf("https://%s", pkg))
		if err == nil {
			return l, nil
		}
	}
	l, err := getLicenseFromLicenseFile(pkg)
	if err == nil {
		return l, nil
	}
	return diligent.License{}, errors.New("failed to find license")
}

func getLicenseFromLicenseFile(pkg string) (diligent.License, error) {
	cmd := exec.Command("go", "get", "-d", fmt.Sprintf("%s/...", pkg))
	err := cmd.Run()
	if err != nil {
		return diligent.License{}, err
	}

	dir := fmt.Sprintf("%s/src/%s", goPath(), pkg)
	filer, err := filer.FromDirectory(dir)
	if err != nil {
		return diligent.License{}, err
	}
	licenses, err := licensedb.Detect(filer)
	if err != nil {
		return diligent.License{}, err
	}
	maxKey := ""
	for k, v := range licenses {
		if maxKey == "" || licenses[maxKey].Confidence < v.Confidence {
			maxKey = k
		}
	}
	if maxKey == "" {
		return diligent.License{}, errors.New("could not identify license")
	}
	return diligent.GetLicenseFromIdentifier(maxKey)
}
