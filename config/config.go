package config

import (
	"fmt"

	"github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"
	"github.com/cloudfoundry/mkman/validators"
)

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	EtcdPath     string   `yaml:"etcd"`
	ConsulPath   string   `yaml:"consul"`
	StubPaths    []string `yaml:"stubs"`
}

func (c Config) Validate() error {
	errors := multierror.NewMultiError("config")
	cfPath := validators.NewValidationTarget(c.CFPath, "cf")
	validator := validators.AllOf(
		validators.NonEmpty(),
		validators.AbsolutePath(),
		validators.AllOf(
			validators.ExistsOnFilesystem(
				validators.Directory(),
			),
		),
	)

	err := cfPath.ValidateWith(validator)
	if err != nil {
		errors.Add(err)
	}

	stemcellPath := validators.NewValidationTarget(c.StemcellPath, "stemcell")
	validator = validators.AllOf(
		validators.NonEmpty(),
		validators.AbsolutePath(),
		validators.AllOf(
			validators.ExistsOnFilesystem(
				validators.File(),
			),
		),
	)

	err = stemcellPath.ValidateWith(validator)
	if err != nil {
		errors.Add(err)
	}

	etcdPath := validators.NewValidationTarget(c.EtcdPath, "etcd")
	validator = validators.AllOf(
		validators.NonEmpty(),
		validators.AnyOf(
			validators.VersionAlias([]string{"director-latest"}),
			validators.AllOf(
				validators.AbsolutePath(),
				validators.ExistsOnFilesystem(
					validators.File(),
					validators.Directory(),
				),
			),
		),
	)

	err = etcdPath.ValidateWith(validator)
	if err != nil {
		errors.Add(err)
	}

	consulPath := validators.NewValidationTarget(c.ConsulPath, "consul")
	err = consulPath.ValidateWith(
		validators.AllOf(
			validators.NonEmpty(),
			validators.AbsolutePath(),
			validators.ExistsOnFilesystem(
				validators.File(),
				validators.Directory(),
			),
		),
	)
	if err != nil {
		errors.Add(err)
	}

	stubErrs := multierror.NewMultiError("stubs")
	stubPaths := validators.NewValidationTarget(c.StubPaths, "stubs")
	emptyErr := validators.NonEmptyArray().Validate(stubPaths)
	if emptyErr != nil {
		errors.Add(emptyErr)
	} else {
		validator = validators.AllOf(
			validators.NonEmpty(),
			validators.AbsolutePath(),
			validators.AllOf(
				validators.ExistsOnFilesystem(
					validators.File(),
				),
			),
		)
		for i, path := range c.StubPaths {
			stubPath := validators.NewValidationTarget(path, fmt.Sprintf("stub at index %d", i))
			err := stubPath.ValidateWith(validator)
			if err != nil {
				stubErrs.Add(err)
			}
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
