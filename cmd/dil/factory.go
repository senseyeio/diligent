package main

import (
	"errors"
	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/govendor"
	"github.com/senseyeio/diligent/npm"
)

var depers = []diligent.Deper{
	npm.New(),
	govendor.New(),
}

func getDeper(filename string, fileContent []byte) (diligent.Deper, error) {
	for _, deper := range depers {
		if deper.IsCompatible(filename, fileContent) {
			return deper, nil
		}
	}
	return nil, errors.New("Unknown file")
}
