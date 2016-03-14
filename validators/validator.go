package validators

import "github.com/cloudfoundry/mkman/Godeps/_workspace/src/github.com/cloudfoundry/multierror"

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

func (vt ValidationTarget) ValidateWith(validator Validator) error {
	err := validator.Validate(vt)
	if err != nil {
		errs := multierror.NewMultiError(vt.name)
		errs.Add(err)
		return errs
	}
	return nil
}

//go:generate counterfeiter . Validator
type Validator interface {
	Validate(vt ValidationTarget) error
	ComposableName() string
}

type FileTypeValidator interface {
	Validator
	FileType() string
}

type Validation struct {
	VersionAliases *[]string
	AllowedType    uint
}
