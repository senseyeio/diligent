package gomod

import (
	"github.com/senseyeio/diligent"
	"github.com/senseyeio/diligent/warning"
	module "github.com/sirkon/goproxy/gomod"
)

type vgo struct {
	lg GoLicenseGetter
}

type GoLicenseGetter interface {
	GetLicense(packagePath string) (diligent.License, error)
}

// New returns a Deper capable of handling dep manifest files
func New(lg GoLicenseGetter) diligent.Deper {
	return &vgo{lg}
}

// Name returns "gomod"
func (v *vgo) Name() string {
	return "gomod"
}

// Dependencies returns the licenses of the go packages defined within the dep manifest
func (v *vgo) Dependencies(file []byte) ([]diligent.Dep, []diligent.Warning, error) {
	mod, err := module.Parse("go.mod", file)
	if err != nil {
		return nil, nil, err
	}

	pkgs := make([]string, 0, len(mod.Require))
	for pkg := range mod.Require {
		pkgs = append(pkgs, pkg)
	}

	for old, new := range mod.Replace {
		for i := range pkgs {
			if old == pkgs[i] {
				switch cast := new.(type) {
				case module.Dependency:
					pkgs[i] = cast.Path
				case module.RelativePath:
					pkgs[i] = string(cast)
				}
			}
		}
	}

	deps := make([]diligent.Dep, 0, len(mod.Require))
	warns := make([]diligent.Warning, 0, len(mod.Require))
	for _, pkg := range pkgs {
		l, err := v.lg.GetLicense(pkg)
		if err != nil {
			warns = append(warns, warning.New(pkg, err.Error()))
		} else {
			deps = append(deps, diligent.Dep{
				Name:    pkg,
				License: l,
			})
		}
	}
	return deps, warns, nil
}

// IsCompatible returns true if the filename is Gopkg.lock
func (v *vgo) IsCompatible(filename string) bool {
	return filename == "go.mod"
}
