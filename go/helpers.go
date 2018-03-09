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
	// look for the base package identifier so github.com/aws/aws-sdk-go/aws becomes github.com/aws/aws-sdk-go
	if len(components) < 3 {
		return diligent.License{}, errors.New("invalid go package path")
	}
	return getLicenseForBasePackage(strings.Join(components[:3], "/"))
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
	return diligent.License{}, errors.New("failed to get license")
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