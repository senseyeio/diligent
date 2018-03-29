package _go

import (
	"github.com/senseyeio/diligent"
	"os/exec"
	"github.com/ryanuber/go-license"
	"fmt"
	"os"
	"github.com/senseyeio/diligent/github"
	"errors"
	"strings"
)

func GetLicense(packagePath string) (diligent.License, error) {
	components := strings.Split(packagePath, "/")
	// in some go vendoring solutions full paths to packages are defined as dependencies
	// need to look for the base package identifier so github.com/aws/aws-sdk-go/aws becomes github.com/aws/aws-sdk-go
	if len(components) < 2 {
		return diligent.License{}, errors.New("invalid go package path")
	}
	// try a three component base package, if possible, as it is most common
	if len(components) >= 3 {
		l, err := getLicenseForBasePackage(strings.Join(components[:3], "/"))
		if err == nil {
			return l, nil
		}
	}
	// can have libraries with just two components, for example gopkg.in/mgo.v2
	return getLicenseForBasePackage(strings.Join(components[:2], "/"))
}

func getLicenseForBasePackage(pkg string) (diligent.License, error) {
	if isGithubPackage(pkg) {
		l, err := getLicenseFromGithub(pkg)
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

func isGithubPackage(pkg string) bool {
	return github.IsGithubURL(fmt.Sprintf("https://%s", pkg))
}

func getLicenseFromGithub(pkg string) (l diligent.License, err error) {
	return github.GetLicenseFromURL(fmt.Sprintf("https://%s", pkg))
}

func getLicenseFromLicenseFile(pkg string) (diligent.License, error) {
	cmd := exec.Command("go", "get", pkg)
	cmd.Run()

	l, err := license.NewFromDir(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkg))
	if err != nil {
		return diligent.License{}, err
	}

	return diligent.GetLicenseFromIdentifier(l.Type)
}