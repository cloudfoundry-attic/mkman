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
		validators.NewPathValidator(validators.Validation{AllowedType: validators.DirType}),
	})

	err := cfPath.Validate()
	if err != nil {
		errors.Add(err)
	}

	return errors
}

// func (c Config) Validate() error {

// 	errors := multierror.NewMultiError("config")

// 	validator := validators.NewValidator(c.CFPath, "cf")
// 	err := validator.Validate(validators.Validation{AllowedType: validators.DirType})
// 	if err != nil {
// 		errors.Add(err)
// 	}

// 	validator = validators.NewValidator(c.StemcellPath, "stemcell")
// 	err = validator.Validate(validators.Validation{AllowedType: validators.FileType})
// 	if err != nil {
// 		errors.Add(err)
// 	}

// 	validator = validators.NewValidator(c.EtcdPath, "etcd")
// 	err = validator.Validate(validators.Validation{
// 		VersionAliases: &[]string{"director-latest"},
// 		AllowedType:    (validators.FileType | validators.DirType),
// 	})
// 	if err != nil {
// 		errors.Add(err)
// 	}

// 	if len(c.StubPaths) < 1 {
// 		errors.Add(fmt.Errorf("value for stubs is required"))
// 	}

// 	stubErrs := multierror.NewMultiError("stubs")
// 	for _, path := range c.StubPaths {
// 		validator = validators.NewValidator(path, path)
// 		err := validator.Validate(validators.Validation{AllowedType: validators.FileType})
// 		if err != nil {
// 			stubErrs.Add(err)
// 		}
// 	}

// 	if stubErrs.Length() > 0 {
// 		errors.Add(stubErrs)
// 	}

// 	if errors.Length() > 0 {
// 		return errors
// 	}
// 	return nil
// }
