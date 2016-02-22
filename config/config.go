package config

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfoundry/mkman/multierror"
)

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	StubPaths    []string `yaml:"stubs"`
}

func (c Config) Validate() multierror.MultiError {
	errors := multierror.MultiError{}

	errors.Add(validatePath(c.CFPath, "cf"))
	errors.Add(validatePath(c.StemcellPath, "stemcell"))

	if len(c.StubPaths) < 1 {
		errors.Add(fmt.Errorf("value for stub path is required"))
	}

	for _, path := range c.StubPaths {
		err := validatePath(path, "stub path")
		errors.Add(err)
		if err.HasAny() {
			break
		}
	}

	return errors
}

func validatePath(object, name string) multierror.MultiError {
	var errors multierror.MultiError
	if object == "" {
		errors.Add(fmt.Errorf("value for %s is required", name))
	}

	if !filepath.IsAbs(object) {
		errors.Add(fmt.Errorf("value for %s must be absolute path", name))
	}

	return errors
}
