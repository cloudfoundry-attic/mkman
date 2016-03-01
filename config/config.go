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

const (
	fileType = "file"
	dirType  = "directory"
)

func (c Config) Validate() error {
	errors := multierror.NewMultiError("config")

	err := validatePath(c.CFPath, "cf", dirType)
	if err != nil {
		errors.Add(err)
	}

	err = validatePath(c.StemcellPath, "stemcell", fileType)
	if err != nil {
		errors.Add(err)
	}

	if len(c.StubPaths) < 1 {
		errors.Add(fmt.Errorf("value for stubs is required"))
	}

	stubErrs := multierror.NewMultiError("stubs")
	for _, path := range c.StubPaths {
		err := validatePath(path, path, fileType)
		if err != nil {
			stubErrs.Add(err)
		}
	}

	if stubErrs.Length() > 0 {
		errors.Add(stubErrs)
	}

	if errors.Length() > 0 {
		return errors
	}
	return nil
}

func validatePath(object, name string, pathType string) error {
	errors := multierror.NewMultiError(name)
	if object == "" {
		errors.Add(fmt.Errorf("value is required"))
		// Return when next error does not make sense
		return errors
	}

	if !filepath.IsAbs(object) {
		errors.Add(fmt.Errorf("value must be absolute path to %s: '%s'", pathType, object))
		// Return when next error does not make sense
		return errors
	}

	stat, err := os.Stat(object)
	if os.IsNotExist(err) {
		errors.Add(fmt.Errorf("%s does not exist: '%s'", pathType, object))
		// Return when next error does not make sense
		return errors
	}

	if stat != nil {
		if stat.IsDir() && pathType == fileType ||
			stat.Mode().IsRegular() && pathType == dirType {
			errors.Add(fmt.Errorf("value must be path to %s: '%s'", pathType, object))
		}
	}

	if errors.Length() > 0 {
		return errors
	}
	return nil
}
