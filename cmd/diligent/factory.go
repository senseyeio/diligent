package main

import (
	"fmt"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/dep"
	"github.com/senseyeio/diligent/github"
	"github.com/senseyeio/diligent/go"
	"github.com/senseyeio/diligent/govendor"
	"github.com/senseyeio/diligent/npm"
)

var (
	gh        = github.New("https://api.github.com")
	goLG      = _go.NewLicenseGetter(gh)
	npmAPIURL = "https://registry.npmjs.org"
)

var depers = []diligent.Deper{
	npm.New(npmAPIURL),
	govendor.New(goLG),
	dep.New(goLG),
}

func getDeper(filename string, fileContent []byte) (diligent.Deper, error) {
	for _, deper := range depers {
		if deper.IsCompatible(filename, fileContent) {
			return deper, nil
		}
	}
	return nil, fmt.Errorf("Diligent does not know how to process '%s' files", filename)
}
