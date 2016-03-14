package validators

import (
	"fmt"
	"strings"
)

type orValidator struct {
	validators []Validator
}

func Or(validators ...Validator) Validator {
	return &orValidator{validators: validators}
}

func (o *orValidator) ComposableName() string {
	var name, delimiter string
	delimiter = " or "
	for _, v := range o.validators {
		name += v.ComposableName() + delimiter
	}
	return strings.TrimSuffix(name, delimiter)
}

func (o *orValidator) Validate(vt ValidationTarget) error {
	for _, v := range o.validators {
		err := v.Validate(vt)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("value must be %s: %s", o.ComposableName(), vt.object)
}
