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
	validator := validators.And(
		validators.NonEmptinessValidator(),
		validators.AbsolutePathValidator(),
		validators.And(
			validators.ExistenceValidator(),
			validators.DirectoryValidator(),
		),
	)

	err := cfPath.ValidateWith(validator)
	if err.Length() > 0 {
		errors.Add(err)
	}

	stemcellPath := validators.NewValidationTarget(c.StemcellPath, "stemcell")
	validator = validators.And(
		validators.NonEmptinessValidator(),
		validators.AbsolutePathValidator(),
		validators.And(
			validators.ExistenceValidator(),
			validators.FileValidator(),
		),
	)

	err = stemcellPath.ValidateWith(validator)
	if err.Length() > 0 {
		errors.Add(err)
	}

	etcdPath := validators.NewValidationTarget(c.EtcdPath, "etcd")
	validator = validators.And(
		validators.NonEmptinessValidator(),
		validators.Or(
			validators.VersionAliasValidator([]string{"director-latest"}),
			validators.And(
				validators.AbsolutePathValidator(),
				validators.ExistenceValidator(),
				validators.Or(
					validators.FileValidator(),
					validators.DirectoryValidator(),
				),
			),
		),
	)

	err = etcdPath.ValidateWith(validator)
	if err.Length() > 0 {
		errors.Add(err)
	}

	stubErrs := multierror.NewMultiError("stubs")
	foo := validators.NewValidationTarget(c.StubPaths, "stubs")
	emptyErr := validators.NonEmptyArrayValidator().Validate(foo)
	if emptyErr != nil {
		errors.Add(emptyErr)
	} else {
		validator = validators.And(
			validators.NonEmptinessValidator(),
			validators.AbsolutePathValidator(),
			validators.And(
				validators.ExistenceValidator(),
				validators.FileValidator(),
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
