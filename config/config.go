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
	EtcdPath     string   `yaml:"etcd"`
	StubPaths    []string `yaml:"stubs"`
}

const (
	none     = 0
	fileType = 1 << iota
	dirType  = 1 << iota
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

	err = validatePath(c.EtcdPath, "etcd", fileType|dirType)
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

func validatePath(object, name string, allowablePathType uint) error {
	translate := func(allowedType uint) string {
		switch allowedType {
		case fileType:
			return "file"
		case dirType:
			return "directory"
		case (fileType | dirType):
			return "file or directory"
		default:
			panic("unhandled")
		}
	}

	errors := multierror.NewMultiError(name)
	if object == "" {
		errors.Add(fmt.Errorf("value is required"))
		// Return when next error does not make sense
		return errors
	}

	if !filepath.IsAbs(object) {
		errors.Add(fmt.Errorf("value must be absolute path to %s: '%s'", translate(allowablePathType), object))
		// Return when next error does not make sense
		return errors
	}

	fileInfo, err := os.Stat(object)
	if os.IsNotExist(err) {
		errors.Add(fmt.Errorf("%s does not exist: '%s'", translate(allowablePathType), object))
		// Return when next error does not make sense
		return errors
	}

	if !isFileTypeAllwed(fileInfo, allowablePathType) {
		errors.Add(fmt.Errorf("value must be path to %s: '%s'", translate(allowablePathType), object))
	}

	if errors.Length() > 0 {
		return errors
	}
	return nil
}

func isFileTypeAllwed(fileInfo os.FileInfo, allowedPathType uint) bool {
	if fileInfo == nil {
		return false
	}

	if fileInfo.Mode().IsRegular() && (allowedPathType&fileType != none) {
		return true
	}

	if fileInfo.Mode().IsDir() && (allowedPathType&dirType != none) {
		return true
	}

	return false
}
