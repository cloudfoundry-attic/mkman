package config

import (
	"fmt"

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
	FileType = 1 << iota
	DirType  = 1 << iota
)

func (c Config) Validate() error {
	errors := multierror.NewMultiError("config")

	validator := NewValidator(c.CFPath, "cf")
	err := validator.Validate(Validation{AllowedType: DirType})
	if err != nil {
		errors.Add(err)
	}

	validator = NewValidator(c.StemcellPath, "stemcell")
	err = validator.Validate(Validation{AllowedType: FileType})
	if err != nil {
		errors.Add(err)
	}

	validator = NewValidator(c.EtcdPath, "etcd")
	err = validator.Validate(Validation{
		VersionAliases: &[]string{"director-latest"},
		AllowedType:    (FileType | DirType),
	})
	if err != nil {
		errors.Add(err)
	}

	if len(c.StubPaths) < 1 {
		errors.Add(fmt.Errorf("value for stubs is required"))
	}

	stubErrs := multierror.NewMultiError("stubs")
	for _, path := range c.StubPaths {
		validator = NewValidator(path, path)
		err := validator.Validate(Validation{AllowedType: FileType})
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
