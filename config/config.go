package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"
)

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	StubPaths    []string `yaml:"stubs"`
}

func (c Config) Validate() error {
	errors := &multierror.MultiError{}

	err := validatePath(c.CFPath, "cf", "directory")
	if err != nil {
		errors.Add(err)
	}

	err = validatePath(c.StemcellPath, "stemcell", "file")
	if err != nil {
		errors.Add(err)
	}

	if len(c.StubPaths) < 1 {
		errors.Add(fmt.Errorf("value for stub path is required"))
	}

	for _, path := range c.StubPaths {
		err := validatePath(path, "stub path", "file")
		if err != nil {
			errors.Add(err)
		}
	}

	if errors.HasAny() {
		return errors
	}
	return nil
}

func validatePath(object, name string, pathType string) error {
	errors := &multierror.MultiError{}
	if object == "" {
		errors.Add(fmt.Errorf("value for %s is required", name))
	}

	if !filepath.IsAbs(object) {
		errors.Add(fmt.Errorf("value for %s must be absolute path to %s: %s", name, pathType, object))
	}

	_, err := os.Stat(object)
	if os.IsNotExist(err) {
		errors.Add(fmt.Errorf("value for %s must be valid path to %s: %s", name, pathType, object))
	}

	if errors.HasAny() {
		return errors
	}
	return nil
}
