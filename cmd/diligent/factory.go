package main

import (
	"fmt"

	"path/filepath"

	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/dep"
	"github.com/senseyeio/diligent/github"
	_go "github.com/senseyeio/diligent/go"
	"github.com/senseyeio/diligent/gomod"
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
	gomod.New(goLG),
}

func getDeper(path string) (diligent.Deper, error) {
	filename := filepath.Base(path)
	for _, deper := range depers {
		if deper.IsCompatible(filename) {
			return deper, nil
		}
	}
	return nil, fmt.Errorf("Diligent does not know how to process '%s' files", filename)
}
