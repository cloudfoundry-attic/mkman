package validators

import "fmt"

type allOf struct {
	validators  []Validator
	firstFailed Validator
}

func AllOf(validators ...Validator) Validator {
	return &allOf{validators: validators}
}

func (a *allOf) ComposableName() string {
	if a.firstFailed == nil {
		return ""
	}
	return a.firstFailed.ComposableName()
}

func (a *allOf) Validate(vt ValidationTarget) error {
	for _, v := range a.validators {
		err := v.Validate(vt)
		if err != nil {
			a.firstFailed = v
			return fmt.Errorf("value must be %s: %s", v.ComposableName(), vt.object)
		}
	}
	return nil
}
