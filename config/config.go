package config

import (
	"github.com/cloudfoundry/mkman/validators"
	"github.com/cloudfoundry/multierror"
)

type Config struct {
	CFPath       string   `yaml:"cf"`
	StemcellPath string   `yaml:"stemcell"`
	EtcdPath     string   `yaml:"etcd"`
	StubPaths    []string `yaml:"stubs"`
}

func (c Config) Validate() error {

	errors := multierror.NewMultiError("config")
	cfPath := validators.NewValidationTarget(c.CFPath, "cf")
	validator := validators.AllOf(
		validators.NonEmpty(),
		validators.AbsolutePath(),
		validators.AllOf(
			validators.ExistsOnFilesystem(),
			validators.Directory(),
		),
	)

	err := cfPath.ValidateWith(validator)
	if err.Length() > 0 {
		errors.Add(err)
	}

	stemcellPath := validators.NewValidationTarget(c.StemcellPath, "stemcell")
	validator = validators.AllOf(
		validators.NonEmpty(),
		validators.AbsolutePath(),
		validators.AllOf(
			validators.ExistsOnFilesystem(),
			validators.File(),
		),
	)

	err = stemcellPath.ValidateWith(validator)
	if err.Length() > 0 {
		errors.Add(err)
	}

	etcdPath := validators.NewValidationTarget(c.EtcdPath, "etcd")
	validator = validators.AllOf(
		validators.NonEmpty(),
		validators.AnyOf(
			validators.VersionAlias([]string{"director-latest"}),
			validators.AllOf(
				validators.AbsolutePath(),
				validators.ExistsOnFilesystem(),
				validators.AnyOf(
					validators.File(),
					validators.Directory(),
				),
			),
		),
	)

	err = etcdPath.ValidateWith(validator)
	if err.Length() > 0 {
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
				validators.ExistsOnFilesystem(),
				validators.File(),
			),
		)
		for _, path := range c.StubPaths {
			stubPath := validators.NewValidationTarget(path, path)
			err := stubPath.ValidateWith(validator)
			if err.Length() > 0 {
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
