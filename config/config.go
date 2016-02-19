package config

import (
	"fmt"

	"github.com/cloudfoundry/mkman/multierror"
)

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	StubPaths    []string `yaml:"stubs"`
}

func (c Config) Validate() *multierror.MultiError {
	errors := multierror.MultiError{}

	if c.CFPath == "" {
		errors.Add(fmt.Errorf("path to cf is missing"))
	}

	if c.StemcellPath == "" {
		errors.Add(fmt.Errorf("path to stemcell is missing"))
	}

	if len(c.StubPaths) < 1 {
		errors.Add(fmt.Errorf("at least one stub path is required"))
	}

	for _, path := range c.StubPaths {
		if path == "" {
			errors.Add(fmt.Errorf("there is an empty stub path"))
			break
		}
	}

	if errors.HasAny() {
		return &errors
	}
	return nil
}
