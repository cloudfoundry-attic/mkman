package validators

import "github.com/cloudfoundry/multierror"

type ValidationTarget struct {
	name       string
	object     interface{}
	validators []Validator
}

func NewValidationTarget(object interface{}, name string) ValidationTarget {
	return ValidationTarget{
		name:       name,
		object:     object,
		validators: []Validator{},
	}
}

func (vt ValidationTarget) Validate() *multierror.MultiError {
	errs := multierror.NewMultiError(vt.name)
	for _, v := range vt.validators {
		err := v.Validate(vt)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs
}

func (vt *ValidationTarget) Add(validators []Validator) {
	for _, v := range validators {
		vt.validators = append(vt.validators, v)
	}
}

type Validator interface {
	Validate(vt ValidationTarget) error
	ComposableName() string
}

type Validation struct {
	VersionAliases *[]string
	AllowedType    uint
}
