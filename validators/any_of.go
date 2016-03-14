package validators

import (
	"fmt"
	"strings"
)

type anyOf struct {
	validators []Validator
}

func AnyOf(validators ...Validator) Validator {
	return &anyOf{validators: validators}
}

func (o *anyOf) ComposableName() string {
	var name, delimiter string
	delimiter = " or "
	for _, v := range o.validators {
		name += v.ComposableName() + delimiter
	}
	return strings.TrimSuffix(name, delimiter)
}

func (o *anyOf) Validate(vt ValidationTarget) error {
	for _, v := range o.validators {
		err := v.Validate(vt)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("value must be %s: %s", o.ComposableName(), vt.object)
}
