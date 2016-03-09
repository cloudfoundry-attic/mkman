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
	cfPath.Add([]validators.Validator{
		validators.NewEmptinessValidator(),
		validators.NewAbsolutePathValidator(),
		validators.NewPathValidator(validators.DirType),
	})

	err := cfPath.Validate()
	if err.Length() > 1 {
		errors.Add(err)
	}

	stemcellPath := validators.NewValidationTarget(c.StemcellPath, "stemcell")
	stemcellPath.Add([]validators.Validator{
		validators.NewEmptinessValidator(),
		validators.NewAbsolutePathValidator(),
		validators.NewPathValidator(validators.FileType),
	})

	err = stemcellPath.Validate()
	if err.Length() > 1 {
		errors.Add(err)
	}

	etcdPath := validators.NewValidationTarget(c.EtcdPath, "etcd")
	etcdPath.Add([]validators.Validator{
		validators.NewEmptinessValidator(),
		validators.And(
			validators.Or(
				validators.NewVersionAliasValidator([]string{"director-latest"}),
				validators.NewAbsolutePathValidator(),
			),
			validators.And(
				validators.NewAbsolutePathValidator(),
				validators.Or(
					validators.NewFilepathValidator(),
					validators.NewDirectoryValidator(),
				),
			),
		),

		// validators.Or(
		// 	validators.And(
		// 		validators.NewAbsolutePathValidator(),
		// 		validators.Or(
		// 			validators.NewFilepathValidator(),
		// 			validators.NewDirectoryValidator(),
		// 		),
		// 	),
		// 	validators.NewVersionAliasValidator([]string{"director-latest"}),
		// ),
	})

	err = etcdPath.Validate()
	if err.Length() > 0 {
		errors.Add(err)
	}

	if errors.Length() > 0 {
		return errors
	}

	return nil
}
