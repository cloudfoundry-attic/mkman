package validators

import "github.com/cloudfoundry/multierror"

type ValidationTarget struct {
	name   string
	object interface{}
}

func NewValidationTarget(object interface{}, name string) ValidationTarget {
	return ValidationTarget{
		name:   name,
		object: object,
	}
}

func (vt ValidationTarget) ValidateWith(validator Validator) *multierror.MultiError {
	errs := multierror.NewMultiError(vt.name)
	err := validator.Validate(vt)
	if err != nil {
		errs.Add(err)
	}
	return errs
}

type Validator interface {
	Validate(vt ValidationTarget) error
	ComposableName() string
}

type Validation struct {
	VersionAliases *[]string
	AllowedType    uint
}
